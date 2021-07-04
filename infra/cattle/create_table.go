package cattle

import (
	"fmt"
	"sort"
	"strings"

	"github.com/8treenet/freedom"
)

type CreateTable struct {
	logger  Logger
	manager *Manager
	name    string
	engine  string
	items   []struct {
		variable string
		kind     string
	}
	order []struct {
		variable string
		num      int
	}
	partitionColumn string
	partitionType   int //1周 2月
}

func (c *CreateTable) init() {
	c.logger = freedom.Logger()
}

func (c *CreateTable) Do() error {
	c.addDefaultColumn()

	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", c.name)
	for i := 0; i < len(c.items); i++ {
		n := ","
		if i == len(c.items)-1 {
			n = ""
		}
		sql += fmt.Sprintf("\t%s %s%s\n", c.items[i].variable, ArrayKind(c.items[i].kind), n)
	}
	sql += fmt.Sprintf(") ENGINE = %s\n", c.engine)

	for c.partitionType != 0 {
		if c.partitionType == 2 {
			sql += fmt.Sprintf("PARTITION BY toYYYYMM(%s)\n", c.partitionColumn)
			break
		}
		sql += fmt.Sprintf("PARTITION BY toMonday(%s)\n", c.partitionColumn)
		break
	}

	sort.Slice(c.order, func(i, j int) bool {
		return c.order[i].num < c.order[j].num
	})

	var orderStrs []string
	for _, v := range c.order {
		orderStrs = append(orderStrs, v.variable)
	}

	if len(orderStrs) > 0 {
		sql += fmt.Sprintf("ORDER BY (%s)", strings.Join(orderStrs, ","))
	}

	c.logger.Infof("CreateTable sql:%s", sql)
	_, err := c.manager.db.Exec(sql)
	return err
}

func (c *CreateTable) AddColumn(variable, kind string, order, partitionType int) {
	c.items = append(c.items, struct {
		variable string
		kind     string
	}{variable, kind})

	if order > 0 {
		c.order = append(c.order, struct {
			variable string
			num      int
		}{variable: variable, num: order})
	}

	if c.partitionColumn != "" || partitionType == 0 {
		return
	}

	if partitionType > 2 && partitionType < 0 {
		panic("partitionType error")
	}
	c.partitionColumn = variable
	c.partitionType = partitionType
}

func (c *CreateTable) addDefaultColumn() {
	c.AddColumn(ColumnRegion, ColumnTypeString, 0, 0)
	c.AddColumn(ColumnCity, ColumnTypeString, 0, 0)
	c.AddColumn(ColumnIP, ColumnTypeIP, 0, 0)
	c.AddColumn(ColumnSourceId, ColumnTypeInt16, 0, 0)

	defPartitionType := 0
	defOrder := 0
	if c.partitionColumn == "" {
		defPartitionType = 2
	}
	if len(c.order) == 0 {
		defOrder = 1
	}
	c.AddColumn(ColumnCreateTime, ColumnTypeDateTime, defOrder, defPartitionType)
}

func (c *CreateTable) SetLogger(l Logger) *CreateTable {
	c.logger = l
	return c
}
