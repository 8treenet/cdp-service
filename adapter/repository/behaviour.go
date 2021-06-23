package repository

import (
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/freedom"
)

func init() {
	behaviourChan = make(chan *entity.Behaviour, 2000)

	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *BehaviourRepository {
			return &BehaviourRepository{}
		})

		initiator.BindBooting(func(bootManager freedom.BootManager) {
			bootManager.RegisterShutdown(func() {
				close(behaviourChan)
				time.Sleep(1 * time.Second)
			})
		})
	})
}

var behaviourChan chan *entity.Behaviour

// BehaviourRepository .
type BehaviourRepository struct {
	freedom.Repository
}

// GetTask .
func (repo *BehaviourRepository) GetTask() <-chan *entity.Behaviour {
	return behaviourChan
}
