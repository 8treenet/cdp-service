package controller

import (
	"time"

	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/freedom"
)

func BehaviourLoop() {
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("behaviourLoop recover:", err)
				BehaviourLoop()
			}
		}()

		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			for {
				cancel := service.BatchProcess()
				if cancel() {
					freedom.Logger().Info("BehaviourLoop cancel")
					break //取消
				}
			}

		})
	}()
}
