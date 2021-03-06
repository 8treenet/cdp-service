package cattle

import (
	"fmt"
	"testing"
	"time"

	"github.com/8treenet/freedom"
	"github.com/go-xorm/builder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDSL(t *testing.T) {
	attributes := map[string]string{}
	tmpl := fmt.Sprintf("%s %s %s", "d", "s", "akk")
	t.Log(attributes, tmpl)
	data := []byte(`<root>
	<and>
	<where column = "haha" compare = "gte"> 20 </where>
	<where column = "haha" compare = "gte"> 50 </where>
	</and>
	</root>`)

	obj, err := NewDSL(data)
	if err != nil {
		panic(err)
	}

	root := obj.node

	walk(root, func(n *node, parent *node) bool {
		fmt.Println(n.XMLName.Local, parent.XMLName.Local, n.Attributes, len(n.Nodes), string(n.Content))
		return true
	})
}

func TestBuilder(t *testing.T) {
	cond1 := builder.And(builder.Eq{"a": "sb"}).And(builder.Eq{"b": 50})
	cond2 := builder.Or(cond1, builder.Eq{"name": "fuck"})
	sql, list, _ := builder.ToSQL(cond2)
	t.Log(sql, list)

	db, e := gorm.Open(mysql.Open("root:123123@tcp(127.0.0.1:3306)/cdp?charset=utf8&parseTime=True&loc=Local&readTimeout=5s&writeTimeout=20s&timeout=2s"), &gorm.Config{})
	if e != nil {
		freedom.Logger().Fatal(e.Error())
	}

	condCdp := builder.Eq{"cdp_source.id": 2}
	sql, err := builder.Select("*").From("cdp_source").Where(condCdp).ToBoundSQL()
	t.Log(sql, err)

	var source struct {
		changes   map[string]interface{}
		ID        int `gorm:"primaryKey;column:id"`
		DDD1      int
		DDD2      string
		Source123 string `gorm:"column:source"`
		DDD3      float32
		Created   time.Time `gorm:"column:created"`
		Updated   time.Time `gorm:"column:updated"`
	}
	err = db.Raw(sql).Scan(&source).Error
	t.Log(source, source.Source123, err)
}

func TestBuilder1(t *testing.T) {
	w := NewWriter()
	builder.Gte{"sb": 1}.WriteTo(w)
	t.Log(builder.ConvertToBoundSQL(w.writer.String(), w.args))
}

func TestCondition(t *testing.T) {
	data := []byte(`<root> 
	<condition>
	<and>
		<where from="user" column = "age" compare = "gte">20</where>
		<or>
			<where from="user" column = "name" compare = "eq">8treenet</where>
			<where from="user" column = "phone" compare = "eq">17897817944</where>
			<and>
				<where from="user" column = "level" compare = "lt">5</where>
				<where from="user" column = "created" compare = "between">2021-06-01 16:11:51,2021-07-13 16:11:51</where>
			</and>
		</or>
	</and>
	</condition>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}
	cond, err := dsl.Condition(dsl.FindConditionNode())
	if err != nil {
		panic(err)
	}
	fmt.Println(builder.ToBoundSQL(cond))
}

func TestSingleSum(t *testing.T) {
	//???????????? ????????????=20 ????????????500 ???????????????
	data := []byte(`<root>
	<from>order</from>
	<join>
		<from leftColumn = "userId" column = "userId">user</from>
		<from leftColumn = "userId" column = "userId">addr</from>
	</join>
	<condition>
		<and>
			<where from="user" column = "age" compare = "gte">20</where>
			<where from="order" column = "price" compare = "gte">500</where>
			<where from="addr" column = "city" compare = "eq">??????</where>
		</and>
	</condition>
	<singleOut column = "price">sum</singleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}
	selectBuilder, err := ExplainSingleAnalysis(dsl, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestSingleOrderCondition(t *testing.T) {
	//????????????????????? ????????????=20 ????????????????????????
	metedata := newMockMetedata("user", func(s string) string {
		switch s {
		case "age":
			return "Int32"
		}
		return ""
	})

	data := []byte(`<root>
	<from>user</from>
	<join>
		<from leftColumn = "userId" column = "userId">order</from>
		<from leftColumn = "userId" column = "userId">addr</from>
	</join>
	<condition>
		<and>
			<where from="user" column = "age" compare = "gte">20</where>
			<where from="order" column = "created" compare = "gt" method = "date">user.created</where>
			<where from="addr" column = "city" compare = "in">??????,??????</where>
		</and>
	</condition>
	<singleOut>people</singleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}
	dsl.SetMetedata(metedata)

	selectBuilder, err := ExplainSingleAnalysis(dsl, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestMultipleSum(t *testing.T) {
	//???????????? ????????????=20 ????????????500 ??????????????? ??????/??????
	data := []byte(`<root>
	<from>order</from>
	<join>
		<from leftColumn = "userId" column = "userId">user</from>
		<from leftColumn = "userId" column = "userId">addr</from>
	</join>
	<condition>
		<and>
			<where from="user" column = "age" compare = "gte">20</where>
			<where from="order" column = "price" compare = "gte">500</where>
			<where from="addr" column = "city" compare = "eq">??????</where>
		</and>
	</condition>
	<multipleOut group = "hour" column = "price">sum</multipleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}
	selectBuilder, err := ExplainMultipleAnalysis(dsl, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestMultipleOrderCondition(t *testing.T) {
	//????????????????????? ????????????=20 ??????????????????????????? ???/??????
	data := []byte(`<root>
	<from>order</from>
	<join>
		<from leftColumn = "userId" column = "userId">user</from>
		<from leftColumn = "userId" column = "userId">addr</from>
	</join>
	<condition>
		<and>
			<where from="user" column = "age" compare = "gte">20</where>
			<where from="order" column = "created" compare = "gt" method = "date">user.created</where>
			<where from="addr" column = "city" compare = "in">??????,??????</where>
		</and>
	</condition>
	<multipleOut group = "productId">count</multipleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}

	selectBuilder, err := ExplainMultipleAnalysis(dsl, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestRegisterCondition(t *testing.T) {
	//??????????????? ??????/??????
	data := []byte(`<root>
	<from>user_register</from>
	<multipleOut group = "sourceId">count</multipleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}

	selectBuilder, err := ExplainMultipleAnalysis(dsl, time.Now().AddDate(0, 0, -500), time.Now().AddDate(0, 0, 100))
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestExplainPersonasAnalysis(t *testing.T) {
	//?????????????????? ?????????????????????1001?????????1002???????????????  3?????? ????????????50??????
	data := []byte(`<list>
	<root>
	<from>order</from>
	<join>
		<from leftColumn = "userId" column = "userId">user</from>
		<from leftColumn = "goods_category_id" column = "category_id">goods</from>
	</join>
	<condition>
		<and>
			<where from="user" column = "age" compare = "gte">20</where>
			<where from="addr" column = "city" compare = "in">??????,??????</where>
			<where from="goods" column = "category_id" compare = "in">1001,1002</where>
		</and>
	</condition>
	<personas>
		<personasOut aggregation = "avg" column = "price" compare = "gte">50</personasOut>
	</personas>
	</root>
	</list>`)
	dsls, err := NewArrayDSL(data)
	if err != nil {
		panic(err)
	}
	dsl := dsls[0]

	selectBuilder, err := ExplainPersonasAnalysis(dsl, []string{"1111", "22222", "33333"}, time.Now().AddDate(0, 0, -10))
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestArrayCondition(t *testing.T) {
	data := []byte(`<root> 
	<from>testing2</from>
	<condition>
		<and>
			<where from="testing2" column = "i32s" compare = "in">1,2</where>
			<where from="testing2" column = "sourceId" compare = "in">1,2,3,4</where>
			<where from="testing2" column = "f1" compare = "between">0.5,0.6</where>
		</and>
	</condition>
	<singleOut>people</singleOut>
	</root>`)
	dsl, err := NewDSL(data)
	if err != nil {
		panic(err)
	}
	t2 := newMockMetedata("testing2", func(s string) string {
		switch s {
		case "i32s":
			return "ArrayInt32"
		case "sourceId":
			return "Int32"
		case "f1":
			return "Float32"
		}
		return ""
	})
	dsl.SetMetedata(t2)

	selectBuilder, err := ExplainSingleAnalysis(dsl, time.Now().AddDate(0, 0, -2), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func newMockMetedata(table string, call func(string) string) *mockMetedata {
	return &mockMetedata{table: table, call: call}
}

type mockMetedata struct {
	table string
	call  func(string) string
}

func (m *mockMetedata) GetWarehouse() string {
	return m.table
}

func (m *mockMetedata) GetColumnType(str string) string {
	return m.call(str)
}
