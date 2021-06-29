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
		Title:     "注册事件",
		Warehouse: "register",
		Metadata:  []vo.ReqFeatureMetadataDTO{},
	}
	body.Metadata = append(body.Metadata, vo.ReqFeatureMetadataDTO{Variable: "mobile", Title: "手机号", Kind: "String"})
	str, _ := req.SetJSONBody(body).ToString()
	t.Log(str)
}

func TestSupportController_GetFeatures(t *testing.T) {
	req := requests.NewHTTPRequest(domain+"/support/features").Get().SetQueryParam("page", 1).SetQueryParam("pageSize", 5)
	str, _ := req.Get().ToString()
	t.Log(str)
}

func TestSupportController_PutFeatureBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/support/feature/2").Put()
	var metadata []vo.ReqFeatureMetadataDTO
	metadata = append(metadata, vo.ReqFeatureMetadataDTO{Variable: "createTime", Title: "注册时间", Kind: "DateTime"})
	metadata = append(metadata, vo.ReqFeatureMetadataDTO{Variable: "ip", Title: "ip地址", Kind: "String"})

	str, _ := req.SetJSONBody(metadata).ToString()
	t.Log(str)
}
