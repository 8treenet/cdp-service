package repository

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"

	"github.com/8treenet/cdp-service/domain/po"
	_ "github.com/8treenet/cdp-service/infra" //Implicit initialization infra
	"github.com/8treenet/cdp-service/utils"
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

	uuid1, _ := utils.GenerateUUID()
	uuid2, _ := utils.GenerateUUID()
	sourceId := 1

	list := []*po.CustomerTemporary{}
	list = append(list, &po.CustomerTemporary{
		UUID:     uuid1,
		UserID:   uuid1 + "user",
		SourceID: sourceId,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
	list = append(list, &po.CustomerTemporary{
		UUID:     uuid2,
		UserID:   uuid2 + "user",
		SourceID: sourceId,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
	err := repo.CreateTempCustomer(list)
	if err != nil {
		panic(err)
	}

	ids := repo.GetExistTempCustomers([]string{uuid1, uuid2}, sourceId)
	if err != nil {
		panic(err)
	}
	t.Log(ids)
}

func TestIP(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *BehaviourRepository
	//获取资源库
	unitTest.FetchRepository(&repo)

	list1 := []string{}
	list2 := []string{}
	list3 := []string{}
	for i := 0; i < 60; i++ {
		list1 = append(list1, fmt.Sprintf("111.117.222.%d", 100+i))
		list2 = append(list2, fmt.Sprintf("113.46.163.%d", 50+i))
		list3 = append(list3, fmt.Sprintf("119.7.146.%d", 35+i))
	}

	for i := 0; i < 50; i++ {
		list555 := []string{}
		list555 = append(list555, list1...)
		list555 = append(list555, list2...)
		list555 = append(list555, list3...)

		postLen := len(list555) - rand.Intn(5)
		res, err := repo.GetIP(list555[0:postLen])
		if err != nil {
			panic(err)
		}
		fmt.Println("fuck1", postLen, len(res))

		for _, v := range res {
			fmt.Println("fuck3", v)
		}
	}
}

func TestFetchBehaviours(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *BehaviourRepository
	//获取资源库
	unitTest.FetchRepository(&repo)
	list, err := repo.FetchBehaviours(1, 1000)
	for i := 0; i < len(list); i++ {
		t.Log(list[i].ID)
	}

	t.Log(list, err)
	//t.Log(repo.TruncateBehaviour())
}

func TestGetFeatureEntit(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *SupportRepository
	unitTest.FetchRepository(&repo)
	t.Log(repo.GetFeatureEntitys())
	t.Log(repo.GetAllFeatureEntity())
}
