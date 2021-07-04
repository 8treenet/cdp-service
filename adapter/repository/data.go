package repository

import (
	"github.com/8treenet/cdp-service/infra/clickhouse"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *DataRepository {
			return &DataRepository{}
		})
	})
}

// DataRepository .
type DataRepository struct {
	freedom.Repository
	Manager *clickhouse.Manager
}

// NewCreateTable
func (repo *DataRepository) NewCreateTable(name string) *clickhouse.CreateTable {
	return repo.Manager.CreateTable(name)
}

// SaveTable
func (repo *DataRepository) SaveTable(cmd *clickhouse.CreateTable) error {
	return cmd.Do()
}

// NewAlterColumn
func (repo *DataRepository) NewAlterColumn(tableName string) *clickhouse.AlterColumn {
	return repo.Manager.AlterColumn(tableName)
}

// SaveColumn
func (repo *DataRepository) SaveColumn(cmd *clickhouse.AlterColumn) error {
	return cmd.Do()
}
