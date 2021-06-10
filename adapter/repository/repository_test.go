package repository

import (
	"testing"

	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"

	_ "github.com/8treenet/cdp-service/infra" //Implicit initialization infra
	"gorm.io/gorm"
)

func getUnitTest() freedom.UnitTest {
	//创建单元测试工具
	unitTest := freedom.NewUnitTest()
	unitTest.InstallDB(func() interface{} {
		db, e := gorm.Open(mysql.Open("root:123123@tcp(127.0.0.1:3306)/template?charset=utf8&parseTime=True&loc=Local"))
		if e != nil {
			freedom.Logger().Fatal(e.Error())
		}
		db = db.Debug()
		return db
	})

	opt := &redis.Options{
		Addr: "127.0.0.1:6379",
	}
	redisClient := redis.NewClient(opt)
	if e := redisClient.Ping().Err(); e != nil {
		freedom.Logger().Fatal(e.Error())
	}
	unitTest.InstallRedis(func() (client redis.Cmdable) {
		return redisClient
	})
	return unitTest
}

// CustomerRepository
func TestAddTempleteField(t *testing.T) {
	//获取单测工具
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *IntermediaryRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	repo.db().Exec("DELETE FROM `cdp_customer_template` WHERE 1").Row()

	e := repo.AddTemplete("name", "String", "", "", 0)
	t.Log(e)
	repo.AddTemplete("sex", "Integer", "dsb", "", 0)
	repo.AddTemplete("age", "Integer", "", "", 0)
	repo.AddTemplete("mobile", "String", "", "", 0)
	repo.AddTemplete("desc", "String", "", "", 0)
	repo.AddTemplete("iq", "Integer", "", "", 0)
	//t.Log(err)
}
