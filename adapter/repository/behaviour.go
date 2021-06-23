package repository

import (
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/freedom"
)

func init() {
	behaviourChan = make(chan *entity.Behaviour, 2000)
	bufferOver = make(chan bool, 1)

	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *BehaviourRepository {
			return &BehaviourRepository{}
		})

		initiator.BindBooting(func(bootManager freedom.BootManager) {
			bootManager.RegisterShutdown(func() {
				//注册程序关闭事件。关闭行为管道
				close(behaviourChan)
				//等全部处理完成的over通知或3秒超时
				select {
				case <-bufferOver:
					break
				case <-time.After(3 * time.Second):
					break
				}

			})
		})
	})
}

var (
	behaviourChan chan *entity.Behaviour
	bufferOver    chan bool
)

// BehaviourRepository .
type BehaviourRepository struct {
	freedom.Repository
}

// FetchBehaviours max:最大数量, duration:等待时间.
func (repo *BehaviourRepository) FetchBehaviours(max int, duration time.Duration) (list []*entity.Behaviour, cancel func() bool) {
	cancel = func() bool { return false }

	for len(list) < max {
		select {
		case obj, ok := <-behaviourChan:
			if !ok {
				cancel = func() bool {
					//返回匿名函数控制over
					bufferOver <- true
					return true
				}
				return
			}

			list = append(list, obj)
		case <-time.After(duration):
			return
		}
	}
	return
}
