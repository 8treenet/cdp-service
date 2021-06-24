package domain

import (
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *BehaviourService {
			return &BehaviourService{fetchTime: 8 * time.Second, fetchCount: 800}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *BehaviourService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// BehaviourService .
type BehaviourService struct {
	Worker              freedom.Worker
	BehaviourRepository *repository.BehaviourRepository
	fetchTime           time.Duration
	fetchCount          int
}

// BatchProcess
func (service *BehaviourService) BatchProcess() func() bool {
	service.Worker.Logger().Debug("BatchProcess")
	list, cancel := service.BehaviourRepository.FetchBehaviours(service.fetchCount, service.fetchTime)
	if len(list) == 0 {
		return cancel
	}

	service.batchProcess(list)
	return cancel
}

// batch
func (service *BehaviourService) batchProcess(list []*entity.Behaviour) {
	service.Worker.Logger().Debug("batch len:", len(list))
}
