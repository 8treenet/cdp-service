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
	featureID, e := b.Worker.IrisContext().URLParamInt("featureId")
	if e != nil {
		return &infra.JSONResponse{Error: e}
	}
	var list []vo.ReqBehaviourDTO
	if e := b.Request.ReadJSON(&list, false); e != nil {
		return &infra.JSONResponse{Error: e}
	}

	if e := b.BehaviourService.CreateBehaviours(featureID, list); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}

//Post handles the Post: /behaviour/csv route.
func (b *BehaviourController) PostCsv() freedom.Result {
	var params struct {
		FeatureId int    `url:"featureId" validate:"required"`
		Key       string `url:"key" validate:"required"`
	}
	if err := b.Request.ReadQuery(&params); err != nil {
		return &infra.JSONResponse{Error: err}
	}

	if e := b.BehaviourService.CreateBehavioursByCSV(params.FeatureId, params.Key); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}
