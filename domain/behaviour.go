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
			return &BehaviourService{}
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
}

// BatchProcess
func (service *BehaviourService) BatchProcess() {
	var list []*entity.Behaviour
	shutdown := false
	for shutdown {
		select {
		case entityObj := <-service.BehaviourRepository.GetTask():
			if entityObj == nil {
				//程序关闭处理
				shutdown = true
				if len(list) > 0 {
					service.batch(list)
				}
				list = make([]*entity.Behaviour, 0)
				continue
			}

			list = append(list, entityObj)
			if len(list) < 500 {
				continue
			}

			service.batch(list)
			list = make([]*entity.Behaviour, 0)
		case <-time.After(5 * time.Second):
			if len(list) == 0 {
				continue
			}

			service.batch(list)
			list = make([]*entity.Behaviour, 0)
		}
	}
}

// batch
func (service *BehaviourService) batch(list []*entity.Behaviour) {

}
