package domain

import (
	"os"
	"strconv"
	"testing"

	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getUnitTest() freedom.UnitTest {
	os.Setenv(freedom.ProfileENV, os.Getenv("GOPATH")+"/src/github.com/8treenet/cdp-service/server/conf")
	//创建单元测试工具
	unitTest := freedom.NewUnitTest()
	unitTest.InstallDB(func() interface{} {
		db, e := gorm.Open(mysql.Open("root:123123@tcp(127.0.0.1:3306)/cdp?charset=utf8&parseTime=True&loc=Local"))
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

func TestTruncate(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *BehaviourService
	unitTest.FetchService(&service)
	service.Truncate()
	//获取资源库
}

func TestExecuteDayJob(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)
	service.ExecuteDayJob()
	//获取资源库
}

func Test123(t *testing.T) {
	fv, err := strconv.ParseFloat("100.09", 64)
	t.Log(fv, err)

	f := 10.54 / 25.123
	t.Log(f)
}
