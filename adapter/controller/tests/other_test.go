package controller_test

import (
	"testing"

	"github.com/8treenet/cdp-service/domain/vo"
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
		Title:     "测试元数据 ",
		Warehouse: "testing2",
		Metadata:  []vo.ReqFeatureMetadataDTO{},
	}

	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "strs", Title: "strs", Kind: "ArrayString"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "f1", Title: "f1", Kind: "Float32"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "i32s", Title: "i32s", Kind: "ArrayInt32"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "ui64s", Title: "ui64s", Kind: "ArrayUInt64"})
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "f64s", Title: "f64s", Kind: "ArrayFloat64"})
	str, _ := req.SetJSONBody(body).ToString()
	t.Log(str)
}

func TestSupportController_PutFeatureBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/feature/10").Put()
	var metadata []vo.ReqFeatureMetadataDTO
	metadata = append(metadata, vo.ReqFeatureMetadataDTO{Variable: "dts", Title: "dts", Kind: "ArrayDateTime", Partition: 0})

	str, _ := req.SetJSONBody(metadata).ToString()
	t.Log(str)
}

func TestSupportController_GetFeatures(t *testing.T) {
	req := requests.NewHTTPRequest(domain+"/support/features").Get().SetQueryParam("page", 1).SetQueryParam("pageSize", 5)
	str, _ := req.Get().ToString()
	t.Log(str)
}
