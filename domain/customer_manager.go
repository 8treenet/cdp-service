package domain

import (
	"github.com/8treenet/crm-service/adapter/repository"
	"github.com/8treenet/crm-service/domain/vo"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *CustomerManagerService {
			return &CustomerManagerService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *CustomerManagerService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// CustomerManagerService .
type CustomerManagerService struct {
	Worker             freedom.Worker
	CustomerRepository *repository.CustomerRepository
}

// GetTempletes 获取客户模板列表.
func (service *CustomerManagerService) GetTempletes() (result []*vo.CustomerTemplate, e error) {
	list, err := service.CustomerRepository.GetTempletes()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(list); i++ {
		result = append(result, vo.NewCustomerTemplateFromPO(list[i]))
	}
	return
}

// AddTempletes 添加客户模板列表.
func (service *CustomerManagerService) AddTempletes(templates []vo.CustomerTemplate) (e error) {
	for _, v := range templates {
		if v.ID != 0 {
			continue
		}
		e = service.CustomerRepository.AddTemplete(v.Name, v.Kind, v.Dict, v.Reg, v.Index, v.Required)
		if e != nil {
			return
		}
	}
	return
}

// UpdateTempleteSort
func (service *CustomerManagerService) UpdateTempleteSort(id int, sort int) (e error) {
	//return service.CustomerRepository.UpdateTempleteSort(id, sort)
	entity, err := service.CustomerRepository.GetTemplete(id)
	if err != nil {
		return err
	}
	entity.SetSort(sort)

	return service.CustomerRepository.SaveTemplete(entity)
}
