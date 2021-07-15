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
	labelRoot        = "root"
	labelAnd         = "and"
	labelOr          = "or"
	labelWhere       = "where"
	labelCondition   = "condition"
	labelFrom        = "from"
	labelJoin        = "join"
	labelSingleOut   = "singleOut"
	labelMultipleOut = "multipleOut"
	labeDenominator  = "denominator"

	attributeColumn     = "column"
	attributeLeftColumn = "leftColumn"
	attributeCompare    = "compare"
	attributeFrom       = "from"
	attributeMethod     = "method"
	attributeGroup      = "group"

	outTypeSumValue    = "sum"
	outTypeAvgValue    = "avg"
	outTypeMaxValue    = "max"
	outTypeMinValue    = "min"
	outTypeCountValue  = "count"
	outTypePeopleValue = "people"

	methodDate = "date"

	AliasPeoples = "peoples"

	groupMin   = "minute"
	groupHour  = "hour"
	groupDay   = "day"
	groupWeek  = "week"
	groupMonth = "month"
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
	walk(dsl.node, func(n *node, parent *node) bool {
		if n.XMLName.Local == labelCondition {
			result = n
			return false
		}
		return true
	})
	return
}

func (dsl *DSL) FindFrom() (result *node) {
	walk(dsl.node, func(n *node, parent *node) bool {
		if parent.XMLName.Local != labelRoot {
			return true //继续
		}
		if n.XMLName.Local == labelFrom {
			result = n
			return false
		}
		return true
	})
	return
}

func (dsl *DSL) FindDenominator() (result *node) {
	walk(dsl.node, func(n *node, parent *node) bool {
		if n.XMLName.Local == labeDenominator {
			result = n
			return false
		}
		return true
	})
	return
}

func (dsl *DSL) FindJoin() (result *node) {
	walk(dsl.node, func(n *node, parent *node) bool {
		if n.XMLName.Local == labelJoin {
			result = n
			return false
		}
		return true
	})
	return
}

func (dsl *DSL) From(selectBuilder *builder.Builder) *builder.Builder {
	return selectBuilder.From(string(dsl.FindFrom().Content))
}

func (dsl *DSL) SingleOut(isPeople *bool) (*builder.Builder, error) {
	*isPeople = false
	fromNode := dsl.FindFrom()
	if fromNode == nil {
		return nil, errors.New("from标签错误")
	}

	singleNode := dsl.FindSingleOut()
	if singleNode == nil {
		return nil, errors.New("single标签错误")
	}

	table := fromNode.GetContent()
	out := singleNode.GetContent()

	column := singleNode.GetAttribute(attributeColumn)
	switch out {
	case outTypePeopleValue:
		*isPeople = true
		return builder.Select(fmt.Sprintf("count(*) as %s", AliasPeoples)), nil
	case outTypeCountValue:
		return builder.Select("count(*) as count"), nil
	case outTypeSumValue:
		return builder.Select(fmt.Sprintf("sum(%s.%s) as sum", table, column)), nil
	case outTypeAvgValue:
		return builder.Select(fmt.Sprintf("avg(%s.%s) as avg", table, column)), nil
	case outTypeMaxValue:
		return builder.Select(fmt.Sprintf("max(%s.%s) as max", table, column)), nil
	case outTypeMinValue:
		return builder.Select(fmt.Sprintf("min(%s.%s) as min", table, column)), nil
	}
	return nil, errors.New("singleOut未知错误")
}

func (dsl *DSL) FindSingleOut() (result *node) {
	for _, v := range dsl.node.Nodes {
		if v.XMLName.Local == labelSingleOut {
			return v
		}
	}
	return
}

func (dsl *DSL) FindMultipleOut() (result *node) {
	for _, v := range dsl.node.Nodes {
		if v.XMLName.Local == labelMultipleOut {
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

func (dsl *DSL) MultipleOut() (result *builder.Builder, e error) {
	fromNode := dsl.FindFrom()
	if fromNode == nil {
		return nil, errors.New("from标签错误")
	}

	multipleNode := dsl.FindMultipleOut()
	if multipleNode == nil {
		return nil, errors.New("multiple标签错误")
	}

	table := fromNode.GetContent()
	out := multipleNode.GetContent()
	column := multipleNode.GetAttribute(attributeColumn)
	group := multipleNode.GetAttribute(attributeGroup)
	var groupExpl []string
	switch group {
	case groupMin:
		selectAlias := fmt.Sprintf("toStartOfMinute(%s.%s) as %s", table, ColumnCreateTime, ColumnCreateTime)
		groupByExpl := fmt.Sprintf("toStartOfMinute(%s.%s)", table, ColumnCreateTime)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	case groupHour:
		selectAlias := fmt.Sprintf("toStartOfHour(%s.%s) as %s", table, ColumnCreateTime, ColumnCreateTime)
		groupByExpl := fmt.Sprintf("toStartOfHour(%s.%s)", table, ColumnCreateTime)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	case groupDay:
		selectAlias := fmt.Sprintf("toStartOfDay(%s.%s) as %s", table, ColumnCreateTime, ColumnCreateTime)
		groupByExpl := fmt.Sprintf("toStartOfDay(%s.%s)", table, ColumnCreateTime)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	case groupWeek:
		selectAlias := fmt.Sprintf("toStartOfWeek(%s.%s) as %s", table, ColumnCreateTime, ColumnCreateTime)
		groupByExpl := fmt.Sprintf("toStartOfWeek(%s.%s)", table, ColumnCreateTime)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	case groupMonth:
		selectAlias := fmt.Sprintf("toStartOfMonth(%s.%s) as %s", table, ColumnCreateTime, ColumnCreateTime)
		groupByExpl := fmt.Sprintf("toStartOfMonth(%s.%s)", table, ColumnCreateTime)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	default:
		selectAlias := fmt.Sprintf("%s.%s as %s", table, group, group)
		groupByExpl := fmt.Sprintf("%s.%s", table, group)
		groupExpl = append(groupExpl, selectAlias, groupByExpl)
	}

	switch out {
	case outTypePeopleValue:
		result = builder.Select(fmt.Sprintf("count(*) as %s, %s", AliasPeoples, groupExpl[0]))
	case outTypeCountValue:
		result = builder.Select(fmt.Sprintf("count(*) as count, %s", groupExpl[0]))
	case outTypeSumValue:
		result = builder.Select(fmt.Sprintf("sum(%s.%s) as sum, %s", table, column, groupExpl[0]))
	case outTypeAvgValue:
		result = builder.Select(fmt.Sprintf("avg(%s.%s) as avg, %s", table, column, groupExpl[0]))
	case outTypeMaxValue:
		result = builder.Select(fmt.Sprintf("max(%s.%s) as max, %s", table, column, groupExpl[0]))
	case outTypeMinValue:
		result = builder.Select(fmt.Sprintf("min(%s.%s) as min, %s", table, column, groupExpl[0]))
	default:
		e = errors.New("multipleOut未知错误")
		return
	}
	if out == outTypePeopleValue {
		result = result.GroupBy(fmt.Sprintf("%s.%s,%s", table, ColumnUserId, groupExpl[1]))
		return
	}
	result = result.GroupBy(groupExpl[1])
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
		value   interface{}
		from    string
		method  string
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
	method = nd.GetAttribute(attributeMethod)
	column = from + "." + column
	value = string(nd.Content)

	switch method {
	case methodDate:
		column = fmt.Sprintf("date(%s)", column)
		value = builder.Expr(fmt.Sprintf("date(%s)", value))
		// value = fmt.Sprintf("date(%s)", value)
	}
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
	case "in":
		list, e := utils.ToInterfaces(strings.Split(value.(string), ","))
		if e != nil {
			err = fmt.Errorf("in error %w", e)
			return
		}
		cond = builder.In(column, list...)
	case "between":
		list := strings.Split(value.(string), ",")
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

func walk(in *node, f func(*node, *node) bool) {
	for _, n := range in.Nodes {
		next := f(n, in)
		if !next {
			break
		}

		if len(n.Nodes) == 0 {
			continue
		}
		walk(n, f)
	}
}
