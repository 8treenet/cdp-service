package controller_test

import (
	"testing"

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
