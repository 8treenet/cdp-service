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

	obj, err := newDSL(data)
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
	dsl, err := newDSL(data)
	if err != nil {
		panic(err)
	}
	cond, err := dsl.Condition(dsl.FindCondition())
	if err != nil {
		panic(err)
	}
	fmt.Println(builder.ToBoundSQL(cond))
}

func TestJoinCondition(t *testing.T) {
	//北京地区 年龄大于=20 订单大于500 的总销售额
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
			<where from="addr" column = "city" compare = "eq">北京</where>
		</and>
	</condition>
	<single column = "price">sum</single>
	</root>`)
	dsl, err := newDSL(data)
	if err != nil {
		panic(err)
	}
	cond, err := dsl.Condition(dsl.FindCondition())
	if err != nil {
		panic(err)
	}

	var isPeople bool
	selectBuilder, err := dsl.SingleSelect(&isPeople)
	if err != nil {
		panic(err)
	}
	if isPeople {
		subSelect := dsl.From(builder.Select(ColumnUserId)).Where(cond).GroupBy(ColumnUserId)
		join, err := dsl.Join(subSelect)
		if err != nil {
			panic(err)
		}
		selectBuilder = selectBuilder.From(join, "people")
	} else {
		join, err := dsl.Join(dsl.From(selectBuilder))
		if err != nil {
			panic(err)
		}
		selectBuilder = join.Where(cond)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestLiucunCondition(t *testing.T) {
	//北京和天津地区 年龄大于=20 的次日下单用户数
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
			<where from="addr" column = "city" compare = "in">北京,天津</where>
		</and>
	</condition>
	<single>people</single>
	</root>`)
	dsl, err := newDSL(data)
	if err != nil {
		panic(err)
	}
	cond, err := dsl.Condition(dsl.FindCondition())
	if err != nil {
		panic(err)
	}

	var isPeople bool
	selectBuilder, err := dsl.SingleSelect(&isPeople)
	if err != nil {
		panic(err)
	}
	if isPeople {
		subSelect := dsl.From(builder.Select(ColumnUserId)).Where(cond).GroupBy(ColumnUserId)
		join, err := dsl.Join(subSelect)
		if err != nil {
			panic(err)
		}
		selectBuilder = selectBuilder.From(join, "people")
	} else {
		join, err := dsl.Join(dsl.From(selectBuilder))
		if err != nil {
			panic(err)
		}
		selectBuilder = join.Where(cond)
	}
	fmt.Println(selectBuilder.ToBoundSQL())
}

func TestDenominator(t *testing.T) {
	//北京地区 年龄大于=20 订单大于500 的总销售额
	data := []byte(`<root>
	<denominator>fuckDDD</denominator>
	</root>`)
	dsl, err := newDSL(data)
	if err != nil {
		panic(err)
	}
	t.Log(dsl.FindDenominator().GetContent())
}
