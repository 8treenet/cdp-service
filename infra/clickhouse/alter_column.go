package clickhouse

import (
	"fmt"
)

type AlterColumn struct {
	manager   *Manager
	tableName string
	variable  string
	kind      string
}

func (ac *AlterColumn) Do() error {
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", ac.tableName, ac.variable, ac.kind)
	_, err := ac.manager.db.Exec(sql)
	return err
}

func (ac *AlterColumn) AddColumn(variable, kind string) {
	ac.variable = variable
	ac.kind = kind
}
