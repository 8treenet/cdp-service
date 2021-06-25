package domain

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
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
}

// 创建渠道 .
func (service *SupportService) CreateSource(source string) error {
	return service.SupportRepo.CreateSouce(source)
}

// GetAllSource .
func (service *SupportService) GetAllSource() ([]*po.Source, error) {
	return service.SupportRepo.GetAllSource()
}
