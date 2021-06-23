package repository

import (
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/freedom"
)

var (
	behaviourChan chan *entity.Behaviour
	bufferOver    chan bool
)

func init() {
	behaviourChan = make(chan *entity.Behaviour, 2000)
	bufferOver = make(chan bool)

	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *BehaviourRepository {
			return &BehaviourRepository{}
		})

		initiator.BindBooting(func(bootManager freedom.BootManager) {

			bootManager.RegisterShutdown(func() {
				//注册程序关闭事件。关闭行为管道，并且等buffer全部处理完成
				close(behaviourChan)
				<-bufferOver
			})
		})
	})
}

// BehaviourRepository .
type BehaviourRepository struct {
	freedom.Repository
}

// ReadStream .
func (repo *BehaviourRepository) ReadStream() <-chan *entity.Behaviour {
	return behaviourChan
}

// TaskOver .
func (repo *BehaviourRepository) TaskOver() {
	repo.Worker().Logger().Debug("TaskOver")
	bufferOver <- true
	return
}
