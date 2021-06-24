package repository

import (
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
	"gorm.io/gorm"
)

func init() {
	behaviourChan = make(chan *po.Behaviour, 3000)
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
	behaviourChan chan *po.Behaviour
	bufferOver    chan bool
)

// BehaviourRepository .
type BehaviourRepository struct {
	freedom.Repository
	GEO *infra.GEO
}

// FetchQueue max:最大数量, duration:等待时间.
func (repo *BehaviourRepository) FetchQueue(max int, duration time.Duration) (list []*po.Behaviour, cancel func() bool) {
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

// AddQueue
func (repo *BehaviourRepository) AddQueue(list []*po.Behaviour) {
	for i := 0; i < len(list); i++ {
		behaviourChan <- list[i]
	}
	return
}

// EnteringWarehouse
func (repo *BehaviourRepository) EnteringWarehouse(list []*po.Behaviour) error {
	return repo.db().CreateInBatches(list, 500).Error
}

// getIP
func (repo *BehaviourRepository) getIP(addr []string) (map[string]*infra.GEOInfo, error) {
	return repo.GEO.ParseBatchIP(addr)
}

// db .
func (repo *BehaviourRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
