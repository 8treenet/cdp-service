package cattle

import (
	"fmt"

	"github.com/8treenet/freedom"
)

type AlterColumn struct {
	logger    Logger
	manager   *Manager
	tableName string
	variable  string
	kind      string
}

func (ac *AlterColumn) init() {
	ac.logger = freedom.Logger()
}

func (ac *AlterColumn) Do() error {
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", ac.tableName, ac.variable, ArrayKind(ac.kind))
	ac.logger.Infof("AlterColumn sql:%s", sql)
	_, err := ac.manager.db.Exec(sql)
	return err
}

func (ac *AlterColumn) AddColumn(variable, kind string) {
	ac.variable = variable
	ac.kind = kind
}

func (ac *AlterColumn) SetLogger(l Logger) {
	ac.logger = l
}
