package cattle

import (
	"database/sql"

	"github.com/8treenet/freedom"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-xorm/builder"
	ormClickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type Logger interface {
	Error(...interface{})
	Errorf(format string, args ...interface{})
	Info(...interface{})
	Infof(format string, args ...interface{})
	Debug(...interface{})
	Debugf(format string, args ...interface{})
}

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		obj := &Manager{}

		initiator.BindInfra(true, obj)
		initiator.InjectController(func(ctx freedom.Context) (com *Manager) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

// Manager
type Manager struct {
	freedom.Infra
	dsn   string
	db    *gorm.DB
	sqlDB *sql.DB
}

// visitConfig .
func (ck *Manager) visitConfig() {
	var cnf struct {
		Addr string `toml:"click_house_addr"`
	}

	if e := freedom.Configure(&cnf, "db.toml"); e != nil {
		freedom.Logger().Fatal("ClickHouse visitConfig", e)
	}
	ck.dsn = cnf.Addr
}

// Booting .
func (ck *Manager) Booting(bootManager freedom.BootManager) {
	ck.visitConfig()
	var e error
	ck.db, e = gorm.Open(ormClickhouse.Open(ck.dsn), &gorm.Config{SkipDefaultTransaction: true})
	if e != nil {
		freedom.Logger().Fatalf("ClickHouse gorm.Open dsn:%s, err:%v", ck.dsn, e.Error())
	}
	ck.sqlDB, e = ck.db.DB()
	if e != nil {
		freedom.Logger().Fatalf("ClickHouse gorm.DB dsn:%s, err:%v", ck.dsn, e.Error())
	}

	freedom.Logger().Debug("ClickHouse connect success dsn:", ck.dsn)
}

func (ck *Manager) CreateTable(name string) *CreateTable {
	reuslt := &CreateTable{manager: ck, name: name, engine: " MergeTree()"}
	reuslt.init()
	return reuslt
}

func (ck *Manager) AlterColumn(tableName string) *AlterColumn {
	reuslt := &AlterColumn{manager: ck, tableName: tableName}
	reuslt.init()
	return reuslt
}

func (ck *Manager) Submit(tableName string) *Submit {
	result := &Submit{tableName: tableName, manager: ck}
	result.init()
	return result
}

func (ck *Manager) CreateQuery(builder *builder.Builder) *Query {
	reuslt := &Query{manager: ck, builder: builder}
	reuslt.init()
	return reuslt
}

func (ck *Manager) tx(f func(*sql.Tx) error) error {
	tx, err := ck.sqlDB.Begin()
	if err != nil {
		return err
	}

	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
