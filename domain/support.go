package domain

import (
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/infra/transaction"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *SupportService {
			return &SupportService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *SupportService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// SupportService 支撑服务 .
type SupportService struct {
	Worker            freedom.Worker
	SupportRepo       *repository.SupportRepository
	FeatureRepository *repository.FeatureRepository
	TX                transaction.Transaction
	DataRepository    *repository.DataRepository
}

// 创建渠道 .
func (service *SupportService) CreateSource(source string) error {
	return service.SupportRepo.CreateSource(source)
}

// GetAllSource .
func (service *SupportService) GetAllSource() ([]*po.Source, error) {
	return service.SupportRepo.GetAllSource()
}

// CreateFeature 创建特征 .
func (service *SupportService) CreateFeature(data vo.ReqFeatureDTO) error {
	entity := service.FeatureRepository.NewFeatureEntity()
	entity.Title = data.Title
	entity.Warehouse = data.Warehouse

	for _, metadata := range data.Metadata {
		entity.FeatureMetadata = append(entity.FeatureMetadata, &po.BehaviourFeatureMetadata{
			Variable:      metadata.Variable,
			Title:         metadata.Title,
			Kind:          metadata.Kind,
			Dict:          metadata.Dict,
			OrderByNumber: metadata.OrderByNumber,
			Partition:     metadata.Partition,
			Created:       time.Now(),
			Updated:       time.Now(),
		})
	}

	return service.TX.Execute(func() error {
		if err := service.FeatureRepository.SaveFeatureEntity(entity); err != nil {
			return err
		}

		cmd := service.DataRepository.NewCreateTable(entity.Warehouse)
		for _, v := range entity.FeatureMetadata {
			cmd.AddColumn(v.Variable, v.Kind, v.OrderByNumber, v.Partition)
		}
		return cmd.Do()
	})
}

// AddFeatureMetadata 为特征添加元数据 .
func (service *SupportService) AddFeatureMetadata(featureId int, list []vo.ReqFeatureMetadataDTO) error {
	entity, err := service.FeatureRepository.GetFeatureEntity(featureId)
	if err != nil {
		return err
	}

	for _, v := range list {
		cmd := service.DataRepository.NewAlterColumn(entity.Warehouse)
		cmd.AddColumn(v.Variable, v.Kind)
		if err := service.DataRepository.SaveColumn(cmd); err != nil {
			service.Worker.Logger().Error(err)
			continue
		}

		entity.AddMetadata(v.Variable, v.Title, v.Kind, v.Dict, 0, 0)
	}

	return service.TX.Execute(func() error {
		return service.FeatureRepository.SaveFeatureEntity(entity)
	})
}

// GetFeaturesByPage 获取Features
func (service *SupportService) GetFeaturesByPage() (result []interface{}, totalPage int, e error) {
	var list []*entity.Feature
	list, totalPage, e = service.FeatureRepository.GetFeatureEntitys()
	for _, v := range list {
		result = append(result, v.View())
	}
	return
}
