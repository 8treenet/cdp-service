package repository

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
	"gorm.io/gorm"
)

func init() {
	behaviourChan = make(chan *po.Behaviour, 5000)
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

// BatchSave
func (repo *BehaviourRepository) BatchSave(list []*po.Behaviour) error {
	return repo.db().CreateInBatches(list, 500).Error
}

// FetchBehaviours
func (repo *BehaviourRepository) FetchBehaviours(featureId int, max int) ([]*entity.Behaviour, error) {
	key := fmt.Sprintf("BehaviourRepository:FetchBehaviours:%d", featureId)
	ok, err := repo.Redis().SetNX(key, 1, time.Minute*1).Result()
	if err != nil || !ok {
		return nil, err
	}
	defer repo.Redis().Del(key)

	list, err := findBehaviourListByWhere(repo, "featureId = ? and processed = ?", []interface{}{featureId, 0}, NewAscLimit("id").SetLength(max))
	if err != nil || len(list) == 0 {
		return nil, err
	}

	var ids []int
	reuslt := []*entity.Behaviour{}
	for i := 0; i < len(list); i++ {
		ids = append(ids, list[i].ID)
		reuslt = append(reuslt, &entity.Behaviour{
			Behaviour: list[i],
		})
	}

	repo.InjectBaseEntitys(reuslt)
	err = repo.db().Model(&po.Behaviour{}).Where("id in (?)", ids).Update("processed", 1).Error
	return reuslt, err
}

// BehaviourSuccess .
func (repo *BehaviourRepository) BehavioursSuccess(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return repo.db().Model(&po.Behaviour{}).Where("id in (?)", ids).Update("processed", 2).Error
}

// BehaviourSuccess .
func (repo *BehaviourRepository) BehavioursError(ids []int) error {
	if len(ids) == 0 {
		return nil
	}
	return repo.db().Model(&po.Behaviour{}).Where("id in (?)", ids).Update("processed", 3).Error
}

// TruncateBehaviour .
func (repo *BehaviourRepository) TruncateBehaviour() bool {
	var count int64
	repo.db().Model(&po.Behaviour{}).Where("processed != 2").Count(&count)
	if count != 0 {
		return false
	}
	sql := fmt.Sprintf("TRUNCATE TABLE %s", (&po.Behaviour{}).TableName())
	if err := repo.db().Exec(sql).Error; err != nil {
		repo.Worker().Logger().Error(err)
		return false
	}
	return true
}

// GetIP
func (repo *BehaviourRepository) GetIP(addr []string) (map[string]*infra.GEOInfo, error) {
	if len(addr) == 0 {
		return make(map[string]*infra.GEOInfo), nil
	}
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
