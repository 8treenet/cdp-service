package controller

import (
	"time"

	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/cdp-service/server/conf"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(i freedom.Initiator) {
		i.BindBooting(func(bootManager freedom.BootManager) {
			behaviourSaveJob()
			behaviourEnteringHouseJob()
			truncateJob()
		})
	})
}

func behaviourSaveJob() {
	//批量入库
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("behaviourSaveJob recover:", err)
				behaviourSaveJob()
			}
		}()

		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			for {
				cancel := service.BatchSave()
				if cancel() {
					freedom.Logger().Info("behaviourSaveJob cancel")
					break //取消
				}
			}

		})
	}()
}

func behaviourEnteringHouseJob() {
	//扫库插入到ck数仓
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("behaviourEnteringHouseJob recover:", err)
				behaviourSaveJob()
			}
		}()

		sec := time.Second * time.Duration(conf.Get().System.JobEnteringHouseSleep)
		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			for {
				service.EnteringHouse()
				time.Sleep(sec)
			}

		})
	}()
}

func truncateJob() {
	getNextTime := func() time.Time {
		addHour := time.Hour * time.Duration(conf.Get().System.JobTruncateHour)
		date := time.Now().Format("2006-01-02")
		currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
		nextTime := currentDateTime.AddDate(0, 0, 1).Add(addHour) //凌晨4点
		return nextTime
	}

	//定时器清理
	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				freedom.Logger().Error("truncateJob recover:", err)
				behaviourSaveJob()
			}
		}()

		nextTime := getNextTime()
		freedom.Logger().Info("truncateJob next time", nextTime)
		for {
			time.Sleep(15 * time.Second)
			if time.Now().Unix() < nextTime.Unix() {
				continue
			}
			freedom.Logger().Info("truncateJob next time", nextTime)

			freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
				service.Truncate()
			})
		}
	}()
}
