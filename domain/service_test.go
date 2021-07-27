package domain

import (
	"os"
	"testing"
	"time"

	"github.com/8treenet/cdp-service/domain/vo"
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

func TestCreateRegdayAnalysis(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)

	var data vo.ReqCreateAnalysis
	data.Name = "reg_day"
	data.Title = "注册用户日统计"
	data.FeatureId = 1
	data.OutType = "singleOut"
	data.DateRange = 7
	data.DateConservation = 1

	data.XmlData = []byte(`<root>
		<from>user_register</from>
		<singleOut>people</singleOut>
	</root>
	`)

	e := service.CreateAnalysis(data)
	t.Log(e)
	time.Sleep(20 * time.Second)
}

func TestCreateSexRegdayAnalysis(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)

	var data vo.ReqCreateAnalysis
	data.Name = "reg_nan_day"
	data.Title = "注册用户性别女日统计"
	data.FeatureId = 1
	data.OutType = "singleOut"
	data.DateRange = 7
	data.DateConservation = 1
	data.DenominatorAnalysisId = 21

	data.XmlData = []byte(`<root>
		<from>user_register</from>
		<condition>
			<and>
				<where from="user_register" column = "gender" compare = "eq">女</where>
			</and>
		</condition>
		<singleOut>people</singleOut>
	</root>
	`)

	e := service.CreateAnalysis(data)
	t.Log(e)
	time.Sleep(20 * time.Second)
}

func TestCreateSourceRegdayAnalysis(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)

	var data vo.ReqCreateAnalysis
	data.Name = "reg_source_day"
	data.Title = "注册用户渠道区间统计"
	data.FeatureId = 1
	data.OutType = "multipleOut"
	data.DateRange = 7
	data.DateConservation = 1

	data.XmlData = []byte(`<root>
		<from>user_register</from>
		<multipleOut group = "sourceId">count</multipleOut>
	</root>
	`)

	e := service.CreateAnalysis(data)
	t.Log(e)
	time.Sleep(10 * time.Second)
}

func TestCreateMinRegdayAnalysis(t *testing.T) {
	//创建注册用户分钟 区间统计
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)

	var data vo.ReqCreateAnalysis
	data.Name = "reg_min_day"
	data.Title = "注册用户分钟 区间统计"
	data.FeatureId = 1
	data.OutType = "multipleOut"
	data.DateRange = 7
	data.DateConservation = 1
	data.DenominatorAnalysisId = 21

	data.XmlData = []byte(`<root>
		<from>user_register</from>
		<multipleOut group = "minute">count</multipleOut>
	</root>
	`)

	e := service.CreateAnalysis(data)
	t.Log(e)
	time.Sleep(10 * time.Second)
}

func TestAnalysisJob(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)
	service.ExecuteDayJob()
	service.ExecuteRefreshJob()
}

func TestQuerySexRegdayAnalysis(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *AnalysisService
	unitTest.FetchService(&service)
	//注册用户统计
	result, e := service.QueryAnalysis(21)
	if e != nil {
		panic(e)
	}
	jsondata, e := result.MarshalJSON()
	t.Log(string(jsondata), e)

	//注册用户性别女的比例
	result2, e2 := service.QueryAnalysis(22)
	if e2 != nil {
		panic(e2)
	}
	jsondata, e2 = result2.MarshalJSON()
	t.Log(string(jsondata), e)

	//注册用户渠道分布
	result3, e3 := service.QueryAnalysis(23)
	if e3 != nil {
		panic(e3)
	}
	jsondata, e2 = result3.MarshalJSON()
	t.Log(string(jsondata), e)

	//注册用户分钟数分布
	result4, e4 := service.QueryAnalysis(24)
	if e3 != nil {
		panic(e4)
	}
	jsondata, e2 = result4.MarshalJSON()
	t.Log(string(jsondata), e)
}

func TestCreatePersona(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *PersonaService
	unitTest.FetchService(&service)

	var req vo.ReqCreatePersona
	req.DateRange = 10
	req.Name = "tuhao"
	req.Title = "土豪"
	req.XmlData = []byte(`<list>
	<root>
	<from>shop_goods</from>
	<condition>
		<and>
			<where from="shop_goods" column = "sourceId" compare = "eq">2</where>
			<where from="shop_goods" column = "tag" compare = "in">801,810</where>
		</and>
	</condition>
	<personas>
		<personasOut aggregation = "sum" column = "price" compare = "gte">10</personasOut>
	</personas>
	</root>
	<root>
	<from>user_register</from>
	<personas>
		<personasOut aggregation = "count"  compare = "eq">1</personasOut>
	</personas>
	</root>
	</list>`)
	t.Log(service.CreatePersona(req))
	time.Sleep(1 * time.Second)
}

func TestPersonaExecuteDayJob(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *PersonaService
	unitTest.FetchService(&service)
	service.ExecuteDayJob()
}

func TestPersonaExecuteRefreshJob(t *testing.T) {
	unitTest := getUnitTest()
	unitTest.Run()

	var service *PersonaService
	unitTest.FetchService(&service)
	service.ExecuteRefreshJob()
}
