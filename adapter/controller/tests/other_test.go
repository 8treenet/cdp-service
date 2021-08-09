package controller_test

import (
	"testing"

	"cdp-service/domain/vo"

	"github.com/8treenet/freedom/infra/requests"
)

func TestSupportController_PostSource(t *testing.T) {
	var data struct {
		Source string `json:"source"`
	}
	data.Source = "ali"

	req := requests.NewHTTPRequest(domain + "/support/source").Post().SetJSONBody(data)
	str, _ := req.ToString()
	t.Log(str)
}

func TestSupportController_GetSource(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/sources").Get()
	str, _ := req.ToString()
	t.Log(str)
}

func TestSupportController_PostFeature(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/feature").Post()
	body := vo.ReqFeatureDTO{
		Title:     "购买商品 ",
		Warehouse: "shop_goods",
		Metadata:  []vo.ReqFeatureMetadataDTO{},
	}

	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "goodId", Title: "商品id", Kind: "UInt32"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "price", Title: "商品价格", Kind: "Float32"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "tag", Title: "商品标签", Kind: "ArrayUInt16"})
	// body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "i32s", Title: "i32s", Kind: "ArrayInt32"})
	// body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "ui64s", Title: "ui64s", Kind: "ArrayUInt64"})
	// body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "f64s", Title: "f64s", Kind: "ArrayFloat64"})
	str, _ := req.SetJSONBody(body).ToString()
	t.Log(str)
}

func TestBehaviourShopGoods(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/behaviour/list?featureId=5").Post()
	var body vo.ReqBehaviourDTO
	//body.CreateTime = "2021-07-28 17:32:47"
	body.IPAddr = "221.204.133.15"
	body.Source = "ali"
	body.UserKey = "yangshu611113513517944333"

	mdata := map[string]interface{}{}
	mdata["goodId"] = 1001
	mdata["price"] = 23.4
	mdata["tag"] = []int{801}

	// mdata["strs"] = []string{"1", "2", "3"}
	// mdata["f1"] = 0.57
	// mdata["i32s"] = []int{1, 2, 3}
	// mdata["ui64s"] = []int{1, 2, 3}
	// mdata["f64s"] = []float64{100.123, 223.555, 3.1415926}
	// mdata["dts"] = []string{"2021-06-30 17:34:59", "2021-06-20 17:34:59"}
	body.Data = mdata

	str, resp := req.SetJSONBody([]vo.ReqBehaviourDTO{body}).ToString()
	t.Log(str, resp)
}

func TestSupportController_PutFeatureBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/feature/10").Put()
	var metadata []vo.ReqFeatureMetadataDTO
	metadata = append(metadata, vo.ReqFeatureMetadataDTO{Variable: "dts", Title: "dts", Kind: "ArrayDateTime"})

	str, _ := req.SetJSONBody(metadata).ToString()
	t.Log(str)
}

func TestSupportController_GetFeatures(t *testing.T) {
	req := requests.NewHTTPRequest(domain+"/support/features").Get().SetQueryParam("page", 1).SetQueryParam("pageSize", 5)
	str, _ := req.Get().ToString()
	t.Log(str)
}
