package controller_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom/infra/requests"
)

var domain = "http://127.0.0.1:8000/cdp-service"

func TestCustomerManagerController_PostList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/list").Post()
	var list []po.CustomerExtensionTemplate
	list = append(list, po.CustomerExtensionTemplate{
		Name: "score",
		Kind: "Int32",
	})
	list = append(list, po.CustomerExtensionTemplate{
		Name: "star",
		Kind: "UInt32",
	})
	list = append(list, po.CustomerExtensionTemplate{
		Name: "addr",
		Kind: "String",
	})
	list = append(list, po.CustomerExtensionTemplate{
		Name:     "level",
		Kind:     "UInt16",
		Required: 1,
	})

	data, resp := req.SetJSONBody(list).ToString()
	t.Log(data, resp.Error)
}

func TestCustomerManagerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/list").Get()
	var body struct {
		Code int                            `json:"code"`
		Msg  string                         `json:"msg"`
		Data []po.CustomerExtensionTemplate `json:"data,omitempty"`
	}
	resp := req.ToJSON(&body)

	t.Log(resp.Error, body.Code, body.Msg)

	str, _ := json.MarshalIndent(body.Data, "   ", "   ")
	t.Log(string(str))
}

func TestCustomerManagerController_PutSort(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/sort").Put()
	str, _ := req.SetQueryParam("id", 53).SetQueryParam("sort", 1002).ToString()
	t.Log(string(str))
}

func TestCustomerController_Post(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Post()
	var data struct {
		po.Customer
		Extension map[string]interface{} `json:"extension"`
		Birthday2 string                 `json:"birthday"`
	}
	data.Name = "yangshu3333"
	data.Gender = "男"
	data.Birthday2 = "1989-05-13"
	data.Email = "4932004@qq.com"
	data.Phone = "135135179333"
	data.UserKey = "yangshu611113513517944333"
	data.WechatUnionID = "10012133333"
	data.Province = "山西"
	data.City = "太原"
	data.Region = "迎泽区"

	data.Extension = make(map[string]interface{})
	data.Extension["score"] = 100
	data.Extension["star"] = 50
	data.Extension["level"] = 20
	data.Extension["addr"] = "11111"

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/855a9b98d30811eb8441804a1460b6f5").Get()
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_PutBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/855a9b98d30811eb8441804a1460b6f5").Put()
	data := map[string]interface{}{
		"birthday": "1991-05-13",
		"gender":   "女",
		"extension": map[string]interface{}{
			"star": 1158,
			"addr": "西城区",
		},
	}

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_PostList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/list").Post()
	datas := make([]struct {
		po.Customer
		Extension map[string]interface{} `json:"extension"`
		Birthday2 string                 `json:"birthday"`
	}, 5)
	for i := 0; i < 5; i++ {
		datas[i].Name = "yangshuList-" + fmt.Sprint(i)
		datas[i].Gender = "男"
		datas[i].Birthday2 = "1989-05-13"
		datas[i].Email = "4932004@qq.com"
		datas[i].Phone = "13513517939" + fmt.Sprint(i)
		datas[i].UserKey = "yangshu611113513517944333" + fmt.Sprint(i)
		datas[i].WechatUnionID = "10012133333" + fmt.Sprint(i)
		datas[i].Province = "山西"
		datas[i].City = "太原"
		datas[i].Region = "迎泽区"
		datas[i].Extension = make(map[string]interface{})
		datas[i].Extension["score"] = 100 + i
		datas[i].Extension["star"] = 50 + i
		datas[i].Extension["level"] = 20 + i
		datas[i].Extension["addr"] = "刺恒小区sss-" + fmt.Sprint(i)
	}

	str, resp := req.SetJSONBody(datas).ToString()
	t.Log(str, resp)
}

func TestCustomerController_DeleteBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Delete()

	req = req.SetQueryParam("userId", []string{"79297552d32511eba16a804a1460b6f5", "792a6ea8d32511eba16a804a1460b6f5", "792ac47ad32511eba16a804a1460b6f5", "792b2730d32511eba16a804a1460b6f5", "792b6d08d32511eba16a804a1460b6f5", "792bb70ed32511eba16a804a1460b6f5", "792c0376d32511eba16a804a1460b6f5", "792c48ccd32511eba16a804a1460b6f5", "792c8f8ad32511eba16a804a1460b6f5", "792cef8ed32511eba16a804a1460b6f5", "792d4100d32511eba16a804a1460b6f5", "792d8ab6d32511eba16a804a1460b6f5", "792dd05cd32511eba16a804a1460b6f5", "792e1620d32511eba16a804a1460b6f5", "792e5ebed32511eba16a804a1460b6f5", "792ea068d32511eba16a804a1460b6f5", "792ee258d32511eba16a804a1460b6f5", "792f25e2d32511eba16a804a1460b6f5", "792f6a7ad32511eba16a804a1460b6f5", "792fb84ad32511eba16a804a1460b6f5"})
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/list").Get()
	str, _ := req.SetQueryParam("page", 1).SetQueryParam("pageSize", 3).ToString()
	t.Log(str)
}

func TestCustomerController_GetKeys(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/key").Get()
	param := []string{"yangshu6111135135179443332", "yangshu6111135135179443334"}

	str, _ := req.SetQueryParam("userKey", param).ToString()
	t.Log(str)
}

func TestCustomerController_GetPhones(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/phone").Get()
	param := []string{"135135179394", "135135179392"}

	str, _ := req.SetQueryParam("phone", param).ToString()
	t.Log(str)
}

func TestCustomerController_GetWechat(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/wechat").Get()
	param := []string{"1001212", "10012122222"}

	str, _ := req.SetQueryParam("unionId", param).ToString()
	t.Log(str)
}
