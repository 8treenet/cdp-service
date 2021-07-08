module github.com/8treenet/cdp-service

go 1.13

//replace github.com/8treenet/freedom => /Users/ysmac/go/src/github.com/8treenet/freedom

require (
	github.com/8treenet/freedom v1.8.10
	github.com/ClickHouse/clickhouse-go v1.4.5
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/iris-contrib/go.uuid v2.0.0+incompatible
	github.com/jmoiron/sqlx v1.3.4
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/kataras/golog v0.1.2
	github.com/kataras/iris/v12 v12.1.8
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	//go.mongodb.org/mongo-driver v1.5.3
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/mysql v1.1.0
	gorm.io/gorm v1.21.10
)
