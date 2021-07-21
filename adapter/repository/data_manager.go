package repository

import (
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/freedom"
	"github.com/go-xorm/builder"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *DataManagerRepository {
			return &DataManagerRepository{}
		})
	})
}

// DataManagerRepository 数据管理.
type DataManagerRepository struct {
	freedom.Repository
	Manager *cattle.Manager
}

// NewCreateTable
func (repo *DataManagerRepository) NewCreateTable(name string) *cattle.CreateTable {
	return repo.Manager.CreateTable(name).SetLogger(repo.Worker().Logger())
}

// SaveTable
func (repo *DataManagerRepository) SaveTable(cmd *cattle.CreateTable) error {
	return cmd.Do()
}

// NewAlterColumn
func (repo *DataManagerRepository) NewAlterColumn(tableName string) *cattle.AlterColumn {
	return repo.Manager.AlterColumn(tableName).SetLogger(repo.Worker().Logger())
}

// SaveColumn
func (repo *DataManagerRepository) SaveColumn(cmd *cattle.AlterColumn) error {
	return cmd.Do()
}

// NewSubmit
func (repo *DataManagerRepository) NewSubmit(tableName string) *cattle.Submit {
	return repo.Manager.Submit(tableName).SetLogger(repo.Worker().Logger())
}

// SaveSubmit
func (repo *DataManagerRepository) SaveSubmit(submit *cattle.Submit) error {
	return submit.Do()
}

// Query
func (repo *DataManagerRepository) Query(result interface{}, b *builder.Builder) error {
	return repo.Manager.CreateQuery(b).SetLogger(repo.Worker().Logger()).Do(result)
}

// GetRepot
func (repo *DataManagerRepository) GetRepot(AnalysisID int) (report *entity.Report, e error) {
	return
}
