//Package controller generated by 'freedom new-project cdp-service'
package controller

import (
	"cdp-service/domain"
	"cdp-service/domain/po"
	"cdp-service/infra"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindController("/customer/metaDataManager", &CustomerManagerController{})
	})
}

// CustomerManagerController .
type CustomerManagerController struct {
	CustomerTempleteService *domain.CustomerManagerService
	Worker                  freedom.Worker
	Request                 *infra.Request
}

//Post handles the Post: /customer/metaDataManager route.
func (c *CustomerManagerController) Post() freedom.Result {
	var list []po.CustomerExtensionMetadata
	if e := c.Request.ReadJSON(&list, true); e != nil {
		return &infra.JSONResponse{Error: e}
	}

	if e := c.CustomerTempleteService.AddMetaData(list); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}

//Get handles the Get: /customer/metaDataManager route.
func (c *CustomerManagerController) Get() freedom.Result {
	data, e := c.CustomerTempleteService.GetMetaData()
	if e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{Object: data}
}

//PutSort handles the put: /customer/metaDataManager/sort route.
func (c *CustomerManagerController) PutSort() freedom.Result {
	var arg struct {
		ID   int `url:"id" validate:"required"`
		Sort int `url:"sort" validate:"required"`
	}
	err := c.Request.ReadQuery(&arg)
	if err != nil {
		return &infra.JSONResponse{Error: err}
	}

	if e := c.CustomerTempleteService.UpdateMetaDataSort(arg.ID, arg.Sort); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}
