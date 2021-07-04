package repository

import (
	"github.com/8treenet/cdp-service/infra/cattle"
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
	Manager *cattle.Manager
}

// NewCreateTable
func (repo *DataRepository) NewCreateTable(name string) *cattle.CreateTable {
	return repo.Manager.CreateTable(name).SetLogger(repo.Worker().Logger())
}

// SaveTable
func (repo *DataRepository) SaveTable(cmd *cattle.CreateTable) error {
	return cmd.Do()
}

// NewAlterColumn
func (repo *DataRepository) NewAlterColumn(tableName string) *cattle.AlterColumn {
	return repo.Manager.AlterColumn(tableName).SetLogger(repo.Worker().Logger())
}

// SaveColumn
func (repo *DataRepository) SaveColumn(cmd *cattle.AlterColumn) error {
	return cmd.Do()
}

// NewSubmit
func (repo *DataRepository) NewSubmit(tableName string) *cattle.Submit {
	return repo.Manager.Submit(tableName).SetLogger(repo.Worker().Logger())
}

// SaveSubmit
func (repo *DataRepository) SaveSubmit(submit *cattle.Submit) error {
	return submit.Do()
}
