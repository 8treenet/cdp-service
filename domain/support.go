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
	Worker      freedom.Worker
	SupportRepo *repository.SupportRepository
	TX          transaction.Transaction
}

// 创建渠道 .
func (service *SupportService) CreateSource(source string) error {
	return service.SupportRepo.CreateSouce(source)
}

// GetAllSource .
func (service *SupportService) GetAllSource() ([]*po.Source, error) {
	return service.SupportRepo.GetAllSource()
}

// CreateFeature 创建特征 .
func (service *SupportService) CreateFeature(data vo.PostFeatureDTO) error {
	entity := service.SupportRepo.NewFeatureEntity()
	entity.Title = data.Title
	entity.Warehouse = data.Warehouse

	for _, metadata := range data.Metadata {
		entity.FeatureMetadata = append(entity.FeatureMetadata, &po.BehaviourFeatureMetadata{
			Variable: metadata.Variable,
			Title:    metadata.Title,
			Kind:     metadata.Kind,
			Dict:     metadata.Dict,
			Created:  time.Now(),
			Updated:  time.Now(),
		})
	}

	return service.TX.Execute(func() error {
		return service.SupportRepo.SaveFeatureEntity(entity)
	})
}

// AddFeatureMetadata 为特征添加元数据 .
func (service *SupportService) AddFeatureMetadata(featureId int, data vo.PostFeatureMetadataDTO) error {
	entity, err := service.SupportRepo.GetFeatureEntity(featureId)
	if err != nil {
		return err
	}
	entity.AddMetadata(data.Variable, data.Title, data.Kind, data.Dict)
	return service.TX.Execute(func() error {
		return service.SupportRepo.SaveFeatureEntity(entity)
	})
}

// GetFeaturesByPage 获取Features
func (service *CustomerService) GetFeaturesByPage() (result []*entity.Feature, totalPage int, e error) {
	result, totalPage, e = service.GetFeaturesByPage()
	return
}
