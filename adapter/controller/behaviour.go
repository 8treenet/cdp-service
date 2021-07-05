package controller

import (
	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindController("/behaviour", &BehaviourController{})
	})
}

// BehaviourController .
type BehaviourController struct {
	BehaviourService *domain.BehaviourService
	Worker           freedom.Worker
	Request          *infra.Request
}

//Post handles the Post: /behaviour/list route.
func (b *BehaviourController) PostList() freedom.Result {
	var list []vo.ReqBehaviourDTO
	if e := b.Request.ReadJSON(&list, false); e != nil {
		return &infra.JSONResponse{Error: e}
	}

	for _, v := range list {
		b.BehaviourService.CreateBehaviour(v)
	}
	return &infra.JSONResponse{}
}
