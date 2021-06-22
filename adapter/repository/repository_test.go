package repository

import (
	"fmt"
	"testing"
	"time"

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

// CustomerRepository
func TestOne(t *testing.T) {
	//获取单测工具
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *CustomerRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	entity1, err1 := repo.GetCustomerByKey("yangshu611113513517944")
	if err1 != nil {
		panic(err1)
	}

	j, _ := entity1.MarshalJSON()
	t.Log("||||||||||" + string(j) + "||||||||||")

	entity2, err2 := repo.GetCustomerByPhone("13513517944")
	if err2 != nil {
		panic(err2)
	}

	j, _ = entity2.MarshalJSON()
	t.Log("||||||||||" + string(j) + "||||||||||")
	//t.Log(err)

	entity3, err3 := repo.GetCustomerByWechat("1001212")
	if err3 != nil {
		panic(err3)
	}

	j, _ = entity3.MarshalJSON()
	t.Log("||||||||||" + string(j) + "||||||||||")

	if entity1.UserID != entity2.UserID || entity1.UserID != entity3.UserID {
		panic("fuck")
	}
}

func TestTemp(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *CustomerRepository
	//获取资源库
	unitTest.FetchRepository(&repo)

	uuid := fmt.Sprint(time.Now().Unix())
	_, err := repo.CreateTempCustomer(uuid)
	if err != nil {
		panic(err)
	}

	entity, err := repo.GetTempCustomer(uuid)
	if err != nil {
		panic(err)
	}
	str, _ := entity.Marshal()
	t.Log(string(str))
}
