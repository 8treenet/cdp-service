package controller_test

import (
	"encoding/json"
	"testing"

	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom/infra/requests"
)

var domain = "http://127.0.0.1:8000/cdp-service"

func TestCustomerManagerController_PostList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/list").Post()
	var list []vo.CustomerExtendTemplate
	list = append(list, vo.CustomerExtendTemplate{
		Name:     "name",
		Kind:     "String",
		Required: 1,
	})
	list = append(list, vo.CustomerExtendTemplate{
		Name: "age",
		Kind: "Integer",
	})
	list = append(list, vo.CustomerExtendTemplate{
		Name: "sex",
		Kind: "Integer",
	})
	list = append(list, vo.CustomerExtendTemplate{
		Name: "mobile",
		Kind: "String",
	})
	list = append(list, vo.CustomerExtendTemplate{
		Name:     "level",
		Kind:     "Integer",
		Required: 1,
	})

	data, resp := req.SetJSONBody(list).ToString()
	t.Log(data, resp.Error)
}

func TestCustomerManagerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/list").Get()
	var body struct {
		Code int                         `json:"code"`
		Msg  string                      `json:"msg"`
		Data []vo.CustomerExtendTemplate `json:"data,omitempty"`
	}
	resp := req.ToJSON(&body)

	t.Log(resp.Error, body.Code, body.Msg)

	str, _ := json.MarshalIndent(body.Data, "   ", "   ")
	t.Log(string(str))
}

func TestCustomerManagerController_PutSort(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/tmplManager/sort").Put()
	str, _ := req.SetQueryParam("id", 35).SetQueryParam("sort", 1002).ToString()
	t.Log(string(str))
}

func TestCustomerController_Post(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Post()
	data := map[string]interface{}{
		"name":   "yangshu1112",
		"age":    1231,
		"sex":    10,
		"mobile": "13513517544",
		"level":  10,
	}

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_PostList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/list").Post()
	var datas []map[string]interface{}
	datas = append(datas, map[string]interface{}{
		"name":   "angshu1234",
		"age":    123,
		"sex":    1,
		"mobile": "12313517144",
		"level":  2,
	})
	datas = append(datas, map[string]interface{}{
		"name":   "syangshu211s",
		"age":    123,
		"sex":    1,
		"mobile": "99413517344",
		"level":  3,
	})

	str, resp := req.SetJSONBody(datas).ToString()
	t.Log(str, resp)
}

func TestCustomerController_PutBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/60bf4c213ac27730dead6e6a").Put()
	data := map[string]interface{}{
		"name":   "yangshu123",
		"age":    31,
		"level":  10,
		"mobile": "13513517944",
	}

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/60bf4c213ac27730dead6e6a").Get()
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetList(t *testing.T) {

}

func TestCustomerController_DeleteBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/60bf506f9fddb4dd6fcdd806").Delete()
	str, resp := req.ToString()
	t.Log(str, resp)
}
