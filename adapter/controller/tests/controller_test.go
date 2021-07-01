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
	req := requests.NewHTTPRequest(domain + "/customer/metaDataManager/list").Post()
	var list []po.CustomerExtensionMetadata
	list = append(list, po.CustomerExtensionMetadata{
		Variable: "score",
		Kind:     "Int32",
		Title:    "积分",
	})
	list = append(list, po.CustomerExtensionMetadata{
		Variable: "star",
		Kind:     "UInt32",
		Title:    "关注",
	})
	list = append(list, po.CustomerExtensionMetadata{
		Variable: "addr",
		Kind:     "String",
		Title:    "地址",
	})
	list = append(list, po.CustomerExtensionMetadata{
		Variable: "level",
		Kind:     "UInt16",
		Required: 1,
		Title:    "级别",
	})

	data, resp := req.SetJSONBody(list).ToString()
	t.Log(data, resp.Error)
}

func TestCustomerManagerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/metaDataManager/list").Get()
	var body struct {
		Code int                            `json:"code"`
		Msg  string                         `json:"msg"`
		Data []po.CustomerExtensionMetadata `json:"data,omitempty"`
	}
	resp := req.ToJSON(&body)

	t.Log(resp.Error, body.Code, body.Msg)

	str, _ := json.MarshalIndent(body.Data, "   ", "   ")
	t.Log(string(str))
}

func TestCustomerManagerController_PutSort(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customer/metaDataManager/sort").Put()
	str, _ := req.SetQueryParam("id", 60).SetQueryParam("sort", 1002).ToString()
	t.Log(string(str))
}

func TestCustomerController_Post(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Post()
	var data struct {
		po.Customer
		Extension        map[string]interface{} `json:"extension"`
		Birthday2        string                 `json:"birthday"`
		Source           string                 `json:"source"`
		IP               string                 `json:"ip"`
		RegisterDateTime string                 `json:"registerDateTime"`
	}
	data.Name = "yangshu3333"
	data.Gender = "男"
	data.Birthday2 = "1989-05-13"
	data.Email = "4932004@qq.com"
	data.Phone = "135135179333"
	data.UserKey = "yangshu611113513517944333"
	data.WechatUnionID = "10012133333"
	data.RegisterDateTime = "2021-06-01 15:15:15"
	//data.City = "太原"
	//data.Region = "山西"
	data.Source = "ali"
	data.IP = "223.20.180.200"

	data.Extension = make(map[string]interface{})
	data.Extension["score"] = 100
	data.Extension["star"] = 50
	data.Extension["level"] = 20
	data.Extension["addr"] = "11111"

	str, resp := req.SetJSONBody(data).ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/466ff946d4e311eb9337804a1460b6f5").Get()
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_PutBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/466ff946d4e311eb9337804a1460b6f5").Put()
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
		Source    string                 `json:"source"`
		IP        string                 `json:"ip"`
	}, 5)
	for i := 0; i < 5; i++ {
		datas[i].Name = "yangshuList-" + fmt.Sprint(i)
		datas[i].Gender = "男"
		datas[i].Birthday2 = "1989-05-13"
		datas[i].Email = "4932004@qq.com"
		datas[i].Phone = "13513517939" + fmt.Sprint(i)
		datas[i].UserKey = "yangshu611113513517944333" + fmt.Sprint(i)
		datas[i].WechatUnionID = "10012133333" + fmt.Sprint(i)
		datas[i].City = "太原"
		datas[i].Region = "山西"
		datas[i].Source = "ali"
		datas[i].Extension = make(map[string]interface{})
		datas[i].Extension["score"] = 100 + i
		datas[i].Extension["star"] = 50 + i
		datas[i].Extension["level"] = 20 + i
		datas[i].Extension["addr"] = "刺恒小区sss-" + fmt.Sprint(i)
	}
	for i := 0; i < 2; i++ {
		datas[i].Source = ""
	}
	datas[0].IP = "125.38.82.23"
	datas[1].IP = "1.192.119.149"

	str, resp := req.SetJSONBody(datas).ToString()
	t.Log(str, resp)
}

func TestCustomerController_DeleteBy(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers").Delete()

	req = req.SetQueryParam("userId", []string{"827680b8d4e311eb9337804a1460b6f5", "8276373ed4e311eb9337804a1460b6f5"})
	str, resp := req.ToString()
	t.Log(str, resp)
}

func TestCustomerController_GetList(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/list").Get()
	str, _ := req.SetQueryParam("page", 1).SetQueryParam("pageSize", 2).ToString()
	t.Log(str)
}

func TestCustomerController_GetKeys(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/key").Get()
	param := []string{"yangshu611113513517944333", "yangshu6111135135179443330"}

	str, _ := req.SetQueryParam("userKey", param).ToString()
	t.Log(str)
}

func TestCustomerController_GetPhones(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/phone").Get()
	param := []string{"135135179333", "135135179390"}

	str, _ := req.SetQueryParam("phone", param).ToString()
	t.Log(str)
}

func TestCustomerController_GetWechat(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/customers/wechat").Get()
	param := []string{"10012133333", "100121333331"}

	str, _ := req.SetQueryParam("unionId", param).ToString()
	t.Log(str)
}
