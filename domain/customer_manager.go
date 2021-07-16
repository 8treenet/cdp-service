package domain

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/infra/transaction"
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
	FeatureRepository  *repository.FeatureRepository
	TX                 transaction.Transaction
	DataRepository     *repository.DataRepository
}

// GetMetaData 获取客户元数据列表.
func (service *CustomerManagerService) GetMetaData() (result []*po.CustomerExtensionMetadata, e error) {
	result, e = service.CustomerRepository.GetMetaData()
	return
}

// AddMetaData 添加客户元数据列表.
func (service *CustomerManagerService) AddMetaData(templates []po.CustomerExtensionMetadata) (e error) {
	//获取注册行为实体
	entity, err := service.FeatureRepository.GetFeatureEntityByWarehouse("user_register")
	if err != nil {
		return err
	}

	for _, v := range templates {
		cmd := service.DataRepository.NewAlterColumn(entity.Warehouse)
		cmd.AddColumn(v.Variable, v.Kind)
		if err := service.DataRepository.SaveColumn(cmd); err != nil {
			service.Worker.Logger().Error(err)
			continue
		}

		entity.AddMetadata(v.Variable, v.Title, v.Kind, v.Dict, 0, 0)
	}

	e = service.TX.Execute(func() error {
		for _, v := range templates {
			if v.ID != 0 {
				continue
			}
			e = service.CustomerRepository.AddMetaData(v)
			if e != nil {
				return e
			}
		}
		return service.FeatureRepository.SaveFeatureEntity(entity)
	})

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
