package controller_test

import (
	"encoding/json"
	"testing"

	"github.com/8treenet/crm-service/domain/vo"
	"github.com/8treenet/freedom/infra/requests"
)

var domain = "http://127.0.0.1:8000/crm-service"

func TestCustomerManagerController_PostList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/manager/list").Post()
	var list []vo.CustomerTemplate
	list = append(list, vo.CustomerTemplate{
		Name:     "name",
		Kind:     "String",
		Required: 1,
		Index:    1,
	})
	list = append(list, vo.CustomerTemplate{
		Name: "age",
		Kind: "Integer",
	})
	list = append(list, vo.CustomerTemplate{
		Name: "sex",
		Kind: "Integer",
	})
	list = append(list, vo.CustomerTemplate{
		Name:  "mobile",
		Kind:  "String",
		Index: 2,
	})
	list = append(list, vo.CustomerTemplate{
		Name:     "level",
		Kind:     "Integer",
		Required: 1,
		Index:    1,
	})

	data, resp := req.SetJSONBody(list).ToString()
	t.Log(data, resp.Error)
}

func TestCustomerManagerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/manager/list").Get()
	var body struct {
		Code int                   `json:"code"`
		Msg  string                `json:"msg"`
		Data []vo.CustomerTemplate `json:"data,omitempty"`
	}
	resp := req.ToJSON(&body)

	t.Log(resp.Error, body.Code, body.Msg)

	str, _ := json.MarshalIndent(body.Data, "   ", "   ")
	t.Log(string(str))
}
