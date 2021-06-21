package controller_test

/*
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
		Kind: "Integer",
	})
	list = append(list, po.CustomerExtensionTemplate{
		Name: "star",
		Kind: "Integer",
	})
	list = append(list, po.CustomerExtensionTemplate{
		Name: "addr",
		Kind: "String",
	})
	list = append(list, po.CustomerExtensionTemplate{
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
	str, _ := req.SetQueryParam("id", 47).SetQueryParam("sort", 1002).ToString()
	t.Log(string(str))
}

func TestCustomerController_Post(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Post()
	var data struct {
		po.Customer
		Extension map[string]interface{} `json:"extension"`
	}
	data.Name = "yangshu6111"
	data.Age = 32
	data.Gender = 1
	data.Extension = make(map[string]interface{})
	data.Extension["score"] = 100
	data.Extension["star"] = 50
	data.Extension["level"] = 20
	data.Extension["addr"] = "11111"

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/2").Get()
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_PutBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/4").Put()
	data := map[string]interface{}{
		"age":    11,
		"gender": 1,
		"extension": map[string]interface{}{
			"star": 1158,
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
	}, 3)
	for i := 0; i < 3; i++ {
		datas[i].Name = "yangshuList-" + fmt.Sprint(i)
		datas[i].Age = 32 + i
		datas[i].Gender = 1
		datas[i].Extension = make(map[string]interface{})
		datas[i].Extension["score"] = 100 + i
		datas[i].Extension["star"] = 50 + i
		datas[i].Extension["level"] = 20 + i
		datas[i].Extension["addr"] = "11111-" + fmt.Sprint(i)
	}

	str, resp := req.SetJSONBody(datas).ToString()
	t.Log(str, resp)
}

func TestCustomerController_DeleteBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Delete()

	req = req.SetQueryParam("id", []int{12, 11, 10, 9})
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetList(t *testing.T) {

}


*/
