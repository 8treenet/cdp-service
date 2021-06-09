//Package controller generated by 'freedom new-project github.com/8treenet/crm-service'
package controller

import (
	"github.com/8treenet/crm-service/domain"
	"github.com/8treenet/crm-service/domain/vo"
	"github.com/8treenet/crm-service/infra"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindController("/customer/tmplManager", &CustomerManagerController{})
	})
}

// CustomerManagerController .
type CustomerManagerController struct {
	CustomerTempleteService *domain.CustomerManagerService
	Worker                  freedom.Worker
	Request                 *infra.Request
}

//Post handles the Post: /customer/tmplManager/list route.
func (c *CustomerManagerController) PostList() freedom.Result {
	var list []vo.CustomerTemplate
	if e := c.Request.ReadJSON(&list, true); e != nil {
		return &infra.JSONResponse{Error: e}
	}

	if e := c.CustomerTempleteService.AddTempletes(list); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}

//Get handles the Get: /customer/tmplManager/list route.
func (c *CustomerManagerController) GetList() freedom.Result {
	data, e := c.CustomerTempleteService.GetTempletes()
	if e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{Object: data}
}

//PutSort handles the put: /customer/tmplManager/sort route.
func (c *CustomerManagerController) PutSort() freedom.Result {
	var arg struct {
		ID   int `url:"id" validate:"required"`
		Sort int `url:"sort" validate:"required"`
	}
	err := c.Request.ReadQuery(&arg)
	if err != nil {
		return &infra.JSONResponse{Error: err}
	}

	if e := c.CustomerTempleteService.UpdateTempleteSort(arg.ID, arg.Sort); e != nil {
		return &infra.JSONResponse{Error: e}
	}
	return &infra.JSONResponse{}
}
