package repository

import (
	"encoding/json"
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

func TestDataTable(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *DataRepository
	unitTest.FetchRepository(&repo)

	cmd := repo.NewCreateTable("testing1")
	cmd.AddColumn("name", "String", 1, 0)
	cmd.AddColumn("create111Time", "DateTime", 2, 2)

	err := repo.SaveTable(cmd)
	t.Log(err)
}

func TestUserRegisterSubmit(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *DataRepository
	unitTest.FetchRepository(&repo)

	cmd := repo.NewSubmit("user_register")
	cmd.AddMetadata("userId", "String")
	cmd.AddMetadata("name", "String")
	cmd.AddMetadata("email", "String")
	cmd.AddMetadata("phone", "String")
	cmd.AddMetadata("gender", "String")
	cmd.AddMetadata("birthday", "Date")

	mdata := map[string]interface{}{}
	mdata["ip"] = "113.46.163.105"
	mdata["city"] = "北京"
	mdata["region"] = "北京"
	mdata["sourceId"] = 1
	mdata["userId"] = "111223123"
	mdata["name"] = "8treenet"
	mdata["email"] = "4932004@qq.com"
	mdata["phone"] = "13513513522"
	mdata["birthday"] = "1989-01-02"
	mdata["gender"] = "男"

	cmd.AddRow(mdata)
	e := repo.SaveSubmit(cmd)
	t.Log(e)
}

func TestArraySubmit(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var repo *DataRepository
	unitTest.FetchRepository(&repo)

	cmd := repo.NewCreateTable("testing2")
	cmd.AddColumn("strs", "ArrayString", 0, 0)
	cmd.AddColumn("f1", "Float32", 0, 0)
	cmd.AddColumn("i32s", "ArrayInt32", 0, 0)
	cmd.AddColumn("ui64s", "ArrayUInt64", 0, 0)
	cmd.AddColumn("f64s", "ArrayFloat64", 0, 0)
	cmd.AddColumn("dts", "ArrayDateTime", 0, 0)

	err := repo.SaveTable(cmd)
	t.Log(err)

	cmdSubmit := repo.NewSubmit("testing2")
	cmdSubmit.AddMetadata("strs", "ArrayString")
	cmdSubmit.AddMetadata("f1", "Float32")
	cmdSubmit.AddMetadata("i32s", "ArrayInt32")
	cmdSubmit.AddMetadata("ui64s", "ArrayUInt64")
	cmdSubmit.AddMetadata("f64s", "ArrayFloat64")
	cmdSubmit.AddMetadata("dts", "ArrayDateTime")

	mdata := map[string]interface{}{}
	mdata["ip"] = "113.46.163.105"
	mdata["city"] = "北京"
	mdata["region"] = "北京"
	mdata["sourceId"] = 1
	mdata["userId"] = "fuckuser"

	mdata["strs"] = []string{"1", "2", "3"}
	mdata["f1"] = 0.57
	mdata["i32s"] = []int{1, 2, 3}
	mdata["ui64s"] = []int{1, 2, 3}
	mdata["f64s"] = []float64{100.123, 223.555, 3.1415926}
	mdata["dts"] = []string{"2021-06-30 17:34:59", "2021-06-20 17:34:59"}

	jsonData, _ := json.Marshal(mdata)
	t.Log("json data", string(jsonData))
	json.Unmarshal(jsonData, &mdata)
	mdata["createTime"] = time.Now()
	cmdSubmit.AddRow(mdata)
	e := repo.SaveSubmit(cmdSubmit)
	t.Log(e)
}
