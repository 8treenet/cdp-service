package controller_test

import (
	"testing"

	"cdp-service/domain/vo"

	"github.com/8treenet/freedom/infra/requests"
)

// data := []byte(`userKey,createTime,ipAddr,source,stringTest,arrayStringTest,int32Test,uint16Test,float32Test,arrayUInt16Test,arrayInt64Test,arrayFloat64Test,dateTest,datetimeTest
// yangshu611113513517944333,2021-06-30 17:32:47,221.204.133.15,ali,freedom good.,controller|service|repository,-50,128,0.57,300|400|50000,-50000|12233|530000,100.123|223.555|3.1415926,2021-06-3,2021-06-30 17:34:59`)

func TestBehaviourController(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/behaviour/list?featureId=8").Post()
	var body vo.ReqBehaviourDTO
	body.CreateTime = "2021-06-30 17:32:47"
	body.IPAddr = "221.204.133.15"
	body.Source = "ali"
	body.UserKey = "yangshu611113513517944333"

	mdata := map[string]interface{}{}
	mdata["stringTest"] = "freedom good."
	mdata["arrayStringTest"] = []string{"controller", "service", "repository"}
	mdata["int32Test"] = -50
	mdata["uint16Test"] = 128
	mdata["float32Test"] = 0.57
	mdata["arrayUInt16Test"] = []int{300, 400, 50000}
	mdata["arrayInt64Test"] = []int{-50000, 12233, 530000}
	mdata["arrayFloat64Test"] = []float64{100.123, 223.555, 3.1415926}

	mdata["dateTest"] = "2021-06-30"
	mdata["datetimeTest"] = "2021-06-30 17:34:59"

	body.Data = mdata

	str, resp := req.SetJSONBody([]vo.ReqBehaviourDTO{body}).ToString()
	t.Log(str, resp)
}

func TestCreateFeature(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/feature").Post()
	var body vo.ReqFeatureDTO
	body.Title = "类型测试"
	body.Warehouse = "test_all_type"
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "stringTest",
		Title:    "字符串",
		Kind:     "String",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "arrayStringTest",
		Title:    "字符串数组",
		Kind:     "ArrayString",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "int32Test",
		Title:    "整形32",
		Kind:     "Int32",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "uint16Test",
		Title:    "无符号整形16",
		Kind:     "UInt16",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "float32Test",
		Title:    "浮点数32",
		Kind:     "Float32",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "arrayUInt16Test",
		Title:    "无符号16整形数组",
		Kind:     "ArrayUInt16",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "arrayInt64Test",
		Title:    "整形64数组",
		Kind:     "ArrayInt64",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "arrayFloat64Test",
		Title:    "浮点数64数组",
		Kind:     "ArrayFloat64",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "dateTest",
		Title:    "日期",
		Kind:     "Date",
	})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{
		Variable: "datetimeTest",
		Title:    "日期时间",
		Kind:     "DateTime",
	})

	str, resp := req.SetJSONBody(body).ToString()
	t.Log(str, resp)
}
