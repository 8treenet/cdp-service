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
			return &BehaviourService{waitBufferTime: 5 * time.Second, waitBufferCount: 500}
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
	waitBufferTime      time.Duration
	waitBufferCount     int
}

// BatchProcess
func (service *BehaviourService) BatchProcess() {
	service.Worker.Logger().Debug("BatchProcess")
	buffer := []*entity.Behaviour{}

	stop := false
	for !stop {
		select {
		case obj, ok := <-service.BehaviourRepository.GetTask():
			if !ok {
				stop = true //close(chan)
				if len(buffer) > 0 {
					service.batch(buffer)
				}
				buffer = []*entity.Behaviour{}
				service.BehaviourRepository.TaskOver() //优雅关闭
				continue
			}

			buffer = append(buffer, obj)
			if len(buffer) < service.waitBufferCount {
				continue
			}

			service.batch(buffer)
			buffer = []*entity.Behaviour{}

		case <-time.After(5 * time.Second):
			if len(buffer) == 0 {
				continue
			}

			service.batch(buffer)
			buffer = []*entity.Behaviour{}
		}
	}
}

// batch
func (service *BehaviourService) batch(list []*entity.Behaviour) {
	service.Worker.Logger().Debug("batch len:", len(list))
}
