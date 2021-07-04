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

// CreateTable max:最大数量, duration:等待时间.
func (repo *DataRepository) CreateTable(name string) *clickhouse.CreateTable {
	return repo.Manager.CreateTable(name)
}

// SaveTable
func (repo *DataRepository) SaveTable(cmd *clickhouse.CreateTable) error {
	return cmd.Do()
}
