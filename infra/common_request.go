package infra

import (
	"strconv"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindInfra(false, func() *CommonRequest {
			return &CommonRequest{}
		})
		initiator.InjectController(func(ctx freedom.Context) (com *CommonRequest) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

// CommonRequest .
type CommonRequest struct {
	freedom.Infra
}

// BeginRequest .
func (req *CommonRequest) BeginRequest(worker freedom.Worker) {
	req.Infra.BeginRequest(worker)

	page, _ := req.Worker().IrisContext().URLParamInt("page")
	pageSize, _ := req.Worker().IrisContext().URLParamInt("pageSize")
	userIdStr := req.Worker().IrisContext().GetHeader("x-user-id")
	userId, _ := strconv.Atoi(userIdStr)

	req.Worker().Store().Set("commonRequest:page", page)
	req.Worker().Store().Set("commonRequest:pageSize", pageSize)
	req.Worker().Store().Set("commonRequest:userId", userId)
}

// GetPage .
func (req *CommonRequest) GetPage() (page int, pageSize int) {
	return req.Worker().Store().GetIntDefault("commonRequest:page", 1), req.Worker().Store().GetIntDefault("commonRequest:pageSize", 10)
}

// GetUserId .
func (req *CommonRequest) GetUserId() int {
	return req.Worker().Store().GetIntDefault("commonRequest:userId", 0)
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

/* Paging common input parameter structure
type PageForm struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}
*/
