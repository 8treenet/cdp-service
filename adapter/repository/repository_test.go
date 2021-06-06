package repository

import (
	"testing"

	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"

	_ "github.com/8treenet/crm-service/infra" //Implicit initialization infra
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

	var repo *CustomerRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	repo.db().Exec("DELETE FROM `cdp_customer_template` WHERE 1").Row()

	e := repo.AddTempleteField("name", "String", "", "", 1)
	t.Log(e)
	repo.AddTempleteField("sex", "Integer", "dsb", "", 0)
	repo.AddTempleteField("age", "Integer", "", "", 0)
	repo.AddTempleteField("mobile", "String", "", "", 2)
	repo.AddTempleteField("desc", "String", "", "", 3)
	repo.AddTempleteField("iq", "Integer", "", "", 0)
	//t.Log(err)
}

// CustomerRepository
func TestCustomerAdd(t *testing.T) {
	//获取单测工具
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *CustomerRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	customer, err := repo.NewCustomer(map[string]interface{}{
		"name":   "yangshu123",
		"sex":    1,
		"age":    31,
		"mobile": "13513517944",
		"desc":   "fuck123123113",
	})
	t.Log(customer, err)
	customer2, err2 := repo.NewCustomer(map[string]interface{}{
		"name":   "qiaojiaojiao1",
		"sex":    2,
		"age":    35,
		"mobile": "13513511111",
		"desc":   "fuck123123113",
	})

	t.Log(customer2, err2)
}

func TestGetSave(t *testing.T) {
	//获取单测工具
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *CustomerRepository
	//获取资源库
	unitTest.FetchRepository(&repo)

	entity, err := repo.GetCustomer("60bb4e4bcd757c893171bb01")
	if err != nil {
		panic(err)
	}

	djson, _ := entity.MarshalJSON()
	t.Log(entity.Templetes, string(djson))

	entity.Update("age", 25)
	entity.Update("iq", 250)
	t.Log(repo.SaveCustomer(entity))
}
