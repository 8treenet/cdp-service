//Package infra generated by 'freedom new-project cdp-service'
package infra

import (
	"encoding/json"
	"strconv"

	"github.com/8treenet/freedom"
	"github.com/kataras/iris/v12/hero"
)

const (
	ERROR   = 7
	SUCCESS = 0
)

// JSONResponse .
type JSONResponse struct {
	Code             int
	Error            error
	Object           interface{}
	DisableLogOutput bool
}

// Dispatch .
func (jrep JSONResponse) Dispatch(ctx freedom.Context) {
	contentType := "application/json"
	var content []byte

	var body struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}
	body.Data = jrep.Object
	body.Code = jrep.Code
	body.Msg = "成功"

	if jrep.Error != nil {
		body.Msg = jrep.Error.Error()
	}
	if jrep.Error != nil && body.Code == 0 {
		body.Code = getErrorCode(freedom.ToWorker(ctx))
	}

	if jrep.Error != nil && body.Code == 0 {
		body.Code = ERROR
	}

	if content, jrep.Error = json.Marshal(body); jrep.Error != nil {
		content = []byte(jrep.Error.Error())
	}

	ctx.Values().Set("code", strconv.Itoa(body.Code))
	if !jrep.DisableLogOutput {
		ctx.Values().Set("response", string(content))
	}

	hero.DispatchCommon(ctx, 200, contentType, content, nil, nil, true)
}

func setErrorCode(work freedom.Worker, code int) {
	work.Store().Set("response::code", code)
}

func getErrorCode(work freedom.Worker) int {
	return work.Store().GetIntDefault("response::code", 0)
}

type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
