package domain

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/po"
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
	CustomerRepository *repository.IntermediaryRepository
}

// GetMetaData 获取客户元数据列表.
func (service *CustomerManagerService) GetMetaData() (result []*po.CustomerExtensionMetadata, e error) {
	result, e = service.CustomerRepository.GetMetaData()
	return
}

// AddMetaData 添加客户元数据列表.
func (service *CustomerManagerService) AddMetaData(templates []po.CustomerExtensionMetadata) (e error) {
	for _, v := range templates {
		if v.ID != 0 {
			continue
		}
		e = service.CustomerRepository.AddMetaData(v)
		if e != nil {
			return
		}
	}
	return
}

// UpdateMetaDataSort
func (service *CustomerManagerService) UpdateMetaDataSort(id int, sort int) (e error) {
	//return service.CustomerRepository.UpdateTempleteSort(id, sort)
	entity, err := service.CustomerRepository.GetOneMetaData(id)
	if err != nil {
		return err
	}
	entity.SetSort(sort)

	return service.CustomerRepository.SaveMetaData(entity)
}
