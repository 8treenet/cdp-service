package cattle

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/8treenet/cdp-service/utils"
	"github.com/go-xorm/builder"
)

const (
	labelAnd            = "and"
	labelOr             = "or"
	labelWhere          = "where"
	labelCondition      = "condition"
	labelFrom           = "from"
	labelJoin           = "join"
	labelSingle         = "single"
	attributeColumn     = "column"
	attributeLeftColumn = "leftColumn"
	attributeCompare    = "compare"
	attributeFrom       = "from"

	singleSumValue    = "sum"
	singleAvgValue    = "avg"
	singleMaxValue    = "max"
	singleMinValue    = "min"
	singleCountValue  = "count"
	singlePeopleValue = "people"
)

//compare eq相等 neq不相等 gt大于 gte大于等于 lt小于  lte小于等于 between范围
type node struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr"`
	Content    []byte     `xml:",innerxml"`
	Nodes      []*node    `xml:",any"`
}

func (n *node) GetAttribute(name string) string {
	for i := 0; i < len(n.Attributes); i++ {
		if n.Attributes[i].Name.Local == name {
			return n.Attributes[i].Value
		}
	}
	return ""
}

func (n *node) GetContent() string {
	return string(n.Content)
}

func newDSL(data []byte) (*DSL, error) {
	var obj node
	err := xml.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}

	return &DSL{node: &obj}, nil
}

type DSL struct {
	node *node
}

func (dsl *DSL) Condition(conditionNode *node) (builder.Cond, error) {
	if len(conditionNode.Nodes) == 0 {
		return nil, errors.New("not found")
	}
	if len(conditionNode.Nodes) > 1 {
		return nil, errors.New("condition 只能为and或or")
	}
	if conditionNode.Nodes[0].XMLName.Local != labelAnd && conditionNode.Nodes[0].XMLName.Local != labelOr {
		return nil, errors.New("condition 只能为and或or")
	}

	return dsl.logicWhere(conditionNode.Nodes[0], conditionNode.Nodes[0].XMLName.Local)
}

func (dsl *DSL) FindCondition() (result *node) {
	walk([]*node{dsl.node}, func(n *node) bool {
		if n.XMLName.Local == labelCondition {
			result = n
			return false
		}
		return true
	})
	return
}

func (dsl *DSL) FindFrom() (result *node) {
	for _, v := range dsl.node.Nodes {
		if v.XMLName.Local == labelFrom {
			return v
		}
	}
	return
}

func (dsl *DSL) FindJoin() (result *node) {
	for _, v := range dsl.node.Nodes {
		if v.XMLName.Local == labelJoin {
			return v
		}
	}
	return
}

func (dsl *DSL) From(selectBuilder *builder.Builder) *builder.Builder {
	return selectBuilder.From(string(dsl.FindFrom().Content))
}

func (dsl *DSL) SingleSelect(isPeople *bool) (*builder.Builder, error) {
	*isPeople = false
	fromNode := dsl.FindFrom()
	if fromNode == nil {
		return nil, errors.New("from标签错误")
	}

	singleNode := dsl.FindSingle()
	if singleNode == nil {
		return nil, errors.New("single标签错误")
	}

	table := fromNode.GetContent()
	single := singleNode.GetContent()

	column := singleNode.GetAttribute(attributeColumn)
	switch single {
	case singlePeopleValue:
		*isPeople = true
		return builder.Select("count(*)"), nil
	case singleCountValue:
		return builder.Select("count(*)"), nil
	case singleSumValue:
		return builder.Select(fmt.Sprintf("sum(%s.%s)", table, column)), nil
	case singleAvgValue:
		return builder.Select(fmt.Sprintf("avg(%s.%s)", table, column)), nil
	case singleMaxValue:
		return builder.Select(fmt.Sprintf("max(%s.%s)", table, column)), nil
	case singleMinValue:
		return builder.Select(fmt.Sprintf("min(%s.%s)", table, column)), nil
	}
	return nil, errors.New("SingleSelect未知错误")
}

func (dsl *DSL) FindSingle() (result *node) {
	for _, v := range dsl.node.Nodes {
		if v.XMLName.Local == labelSingle {
			return v
		}
	}
	return
}
func (dsl *DSL) Join(fromBuilder *builder.Builder) (result *builder.Builder, err error) {
	result = fromBuilder
	joinNode := dsl.FindJoin()
	if joinNode == nil {
		return
	}
	for _, node := range joinNode.Nodes {
		table, cond, e := dsl.join(fromBuilder.TableName(), node)
		if e != nil {
			err = e
			return
		}

		result = result.LeftJoin(table, cond)
	}
	return
}

func (dsl *DSL) logicWhere(andNode *node, label string) (builder.Cond, error) {
	var cond builder.Cond
	for _, vNode := range andNode.Nodes {
		if !utils.InSlice([]string{labelAnd, labelOr, labelWhere}, vNode.XMLName.Local) {
			return nil, errors.New("未识别的标签:" + vNode.XMLName.Local)
		}

		var fcond builder.Cond
		if vNode.XMLName.Local == labelWhere {
			wnode, err := dsl.where(vNode)
			if err != nil {
				return nil, err
			}
			fcond = wnode
		} else {
			lnode, err := dsl.logicWhere(vNode, vNode.XMLName.Local)
			if err != nil {
				return nil, err
			}
			fcond = lnode
		}

		if cond == nil {
			cond = fcond
			continue
		}

		if label == labelAnd {
			cond = cond.And(fcond)
			continue
		}
		cond = cond.Or(fcond)
	}

	if cond == nil {
		return nil, errors.New("not found")
	}
	return cond, nil
}

func (dsl *DSL) where(nd *node) (cond builder.Cond, err error) {
	var (
		column  string
		compare string
		value   string
		from    string
	)

	//后续可以在这里做元数据检查
	attributes := map[string]string{}
	for i := 0; i < len(nd.Attributes); i++ {
		name := nd.Attributes[i].Name.Local
		attributes[name] = nd.Attributes[i].Value
	}
	if column = nd.GetAttribute(attributeColumn); column == "" {
		err = errors.New("attributeColumn 不存在")
		return
	}

	if compare = nd.GetAttribute(attributeCompare); compare == "" {
		err = errors.New("attributecompare 不存在")
		return
	}
	if from = nd.GetAttribute(attributeFrom); from == "" {
		err = errors.New("from 不存在")
		return
	}
	column = from + "." + column

	value = string(nd.Content)
	switch compare {
	case "eq":
		cond = builder.Eq{column: value}
	case "neq":
		cond = builder.Neq{column: value}
	case "gt":
		cond = builder.Gt{column: value}
	case "gte":
		cond = builder.Gte{column: value}
	case "lt":
		cond = builder.Lt{column: value}
	case "lte":
		cond = builder.Lte{column: value}
	case "between":
		list := strings.Split(value, ",")
		if len(list) != 2 {
			err = errors.New("between错误")
			return
		}
		cond = builder.Between{Col: column, LessVal: list[0], MoreVal: list[1]}
	}
	return
}

func (dsl *DSL) join(tableName string, joinFrom *node) (table, cond string, err error) {
	var (
		column     string
		leftColumn string
	)

	table = joinFrom.GetContent()
	if joinFrom.XMLName.Local != labelFrom {
		err = errors.New("错误的join标签 :" + joinFrom.XMLName.Local)
		return
	}

	attributes := map[string]string{}
	for i := 0; i < len(joinFrom.Attributes); i++ {
		name := joinFrom.Attributes[i].Name.Local
		attributes[name] = joinFrom.Attributes[i].Value
	}

	if column = joinFrom.GetAttribute(attributeColumn); column == "" {
		err = errors.New("attributeColumn 不存在")
		return
	}

	if leftColumn = joinFrom.GetAttribute(attributeLeftColumn); leftColumn == "" {
		err = errors.New("attributeLeftColumn 不存在")
		return
	}
	cond = fmt.Sprintf("%s.%s = %s.%s", tableName, leftColumn, table, column)
	return
}

func walk(nodes []*node, f func(*node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}
