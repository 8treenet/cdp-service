package controller

import (
	"time"

	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(i freedom.Initiator) {
		i.BindBooting(func(bootManager freedom.BootManager) {
			behaviourSaveJob()
			behaviourEnteringJob()
		})
	})
}

func behaviourSaveJob() {
	//批量入库
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("behaviourLoop recover:", err)
				behaviourSaveJob()
			}
		}()

		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			for {
				cancel := service.BatchSave()
				if cancel() {
					freedom.Logger().Info("BehaviourLoop cancel")
					break //取消
				}
			}

		})
	}()
}

func behaviourEnteringJob() {
	//扫库插入到ck数仓
}
