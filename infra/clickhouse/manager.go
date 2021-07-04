package clickhouse

import (
	"context"
	"time"

	"github.com/8treenet/freedom"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

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
	dsn string
	db  *sqlx.DB
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
	connect, err := sqlx.Open("clickhouse", ck.dsn)
	if err != nil {
		freedom.Logger().Fatalf("ClickHouse dsn:%s, err:%v", ck.dsn, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if e := connect.PingContext(ctx); e != nil {
		freedom.Logger().Fatalf("ClickHouse ping dsn:%s, err:%v", ck.dsn, e.Error())
	}

	freedom.Logger().Debug("ClickHouse connect success dsn:", ck.dsn)
	ck.db = connect
}

func (ck *Manager) CreateTable(name string) *CreateTable {
	return &CreateTable{manager: ck, name: name, engine: " MergeTree()"}
}

func (ck *Manager) AlterColumn(tableName string) *AlterColumn {
	return &AlterColumn{manager: ck, tableName: tableName}
}
