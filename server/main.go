//Package main generated by 'freedom new-project github.com/8treenet/cdp-service'
package main

import (
	"time"

	_ "github.com/8treenet/cdp-service/adapter/controller" //Implicit initialization controller
	_ "github.com/8treenet/cdp-service/adapter/repository" //Implicit initialization repository
	_ "github.com/8treenet/cdp-service/infra"              //Implicit initialization infra
	localMiddleware "github.com/8treenet/cdp-service/middleware"

	"github.com/8treenet/cdp-service/server/conf"
	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/infra/requests"
	"github.com/8treenet/freedom/middleware"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	app := freedom.NewApplication()
	installMiddleware(app)
	runner := app.NewH2CRunner(conf.Get().App.Other["listen_addr"].(string))
	app.InstallParty("/cdp-service")
	liveness(app)
	installDatabase(app)
	installRedis(app)
	app.Run(runner, *conf.Get().App)
}

func installMiddleware(app freedom.Application) {
	app.InstallMiddleware(middleware.NewRecover())
	app.InstallMiddleware(middleware.NewTrace("x-request-id"))
	//One Loger per request New.
	app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))
	//The middleware output of the log line.
	app.Logger().Handle(localMiddleware.NewLogrusMiddleware(conf.Get().App.Other["logger_path"].(string), conf.Get().App.Other["logger_console"].(bool)))

	//Install the Prometheus middleware.
	middle := middleware.NewClientPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
	requests.InstallMiddleware(middle)

	//HTTP request link middleware that controls the header transmission of requests.
	app.InstallBusMiddleware(middleware.NewBusFilter())
}

func installDatabase(app freedom.Application) {
	app.InstallDB(func() interface{} {
		conf := conf.Get().DB
		db, e := gorm.Open(mysql.Open(conf.Addr), &gorm.Config{})
		if e != nil {
			freedom.Logger().Fatal(e.Error())
		}

		sqlDB, err := db.DB()
		if err != nil {
			freedom.Logger().Fatal(e.Error())
		}
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifeTime) * time.Second)
		return db
	})
}

func installRedis(app freedom.Application) {
	app.InstallRedis(func() (client redis.Cmdable) {
		cfg := conf.Get().Redis
		opt := &redis.Options{
			Addr:               cfg.Addr,
			Password:           cfg.Password,
			DB:                 cfg.DB,
			MaxRetries:         cfg.MaxRetries,
			PoolSize:           cfg.PoolSize,
			ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
			IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
			IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
			MaxConnAge:         time.Duration(cfg.MaxConnAge) * time.Second,
			PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
		}
		redisClient := redis.NewClient(opt)
		if e := redisClient.Ping().Err(); e != nil {
			freedom.Logger().Fatal(e.Error())
		}
		client = redisClient
		return
	})
}

func liveness(app freedom.Application) {
	app.Iris().Get("/ping", func(ctx freedom.Context) {
		ctx.WriteString("pong")
	})
}
