package controller_test

import (
	"testing"

	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom/infra/requests"
)

func TestBehaviourController(t *testing.T) {
	req := requests.NewHTTPRequest(domain + "/behaviour/list").Post()
	var body vo.ReqBehaviourDTO
	body.CreateTime = "2021-06-30 17:32:47"
	body.FeatureID = 10
	body.IPAddr = "221.204.133.15"
	body.Source = "ali"
	body.UserKey = "yangshu611113513517944333"

	mdata := map[string]interface{}{}
	mdata["strs"] = []string{"1", "2", "3"}
	mdata["f1"] = 0.57
	mdata["i32s"] = []int{1, 2, 3}
	mdata["ui64s"] = []int{1, 2, 3}
	mdata["f64s"] = []float64{100.123, 223.555, 3.1415926}
	mdata["dts"] = []string{"2021-06-30 17:34:59", "2021-06-20 17:34:59"}
	body.Data = mdata

	str, resp := req.SetJSONBody([]vo.ReqBehaviourDTO{body}).ToString()
	t.Log(str, resp)
}
