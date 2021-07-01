package controller

import (
	"time"

	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(i freedom.Initiator) {
		i.BindBooting(func(bootManager freedom.BootManager) {
			behaviourJob()
		})
	})
}
func behaviourJob() {
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("behaviourLoop recover:", err)
				behaviourJob()
			}
		}()

		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			for {
				cancel := service.EnteringWarehouse()
				if cancel() {
					freedom.Logger().Info("BehaviourLoop cancel")
					break //取消
				}
			}

		})
	}()
}
