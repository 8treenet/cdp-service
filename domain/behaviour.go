package domain

import (
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *BehaviourService {
			return &BehaviourService{fetchTime: 3 * time.Second, fetchCount: 800}
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

// EnteringWarehouse //批量入库
func (service *BehaviourService) EnteringWarehouse() func() bool {
	service.Worker.Logger().Debug("BatchProcess")
	list, cancel := service.BehaviourRepository.FetchQueue(service.fetchCount, service.fetchTime)
	if len(list) == 0 {
		return cancel
	}

	for i := 0; i < 2; i++ {
		err := service.BehaviourRepository.EnteringWarehouse(list)
		if err == nil {
			break
		}
		if err != nil && i == 1 {
			service.Worker.Logger().Error(err)
		}
		time.Sleep(3 * time.Second)
	}
	return cancel
}
