package controller

import (
	"time"

	"github.com/8treenet/cdp-service/domain"
	"github.com/8treenet/cdp-service/server/conf"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(i freedom.Initiator) {
		i.BindBooting(func(bootManager freedom.BootManager) {
			behaviourSaveJob()
			behaviourEnteringHouseJob()
			truncateJob()
			analysisDayJob()
			analysisRefreshJob()
		})
	})
}

func behaviourSaveJob() {
	//批量入库
	go func() {
		time.Sleep(1 * time.Second)

		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("behaviourSaveJob recover:%v \n%s", err, strace)
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
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("behaviourEnteringHouseJob recover:%v \n%s", err, strace)
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
	if conf.Get().System.JobTruncateHour == 0 {
		return
	}

	//定时器清理
	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("truncateJob recover:%v \n%s", err, strace)
				behaviourSaveJob()
			}
		}()
		dayTrigger(func() {
			freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
				service.Truncate()
			})
		}, conf.Get().System.JobTruncateHour, "truncateJob")
	}()
}

func analysisDayJob() {
	if conf.Get().System.JobAnalysisHour == 0 {
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("analysisDayJob recover:%v \n%s", err, strace)
				analysisDayJob()
			}
		}()
		dayTrigger(func() {
			freedom.ServiceLocator().Call(func(service *domain.AnalysisService) {
				service.ExecuteDayJob()
			})
		}, conf.Get().System.JobAnalysisHour, "analysisDayJob")
	}()
}

func analysisRefreshJob() {
	if conf.Get().System.JobAnalysisRefreshHour == 0 {
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("analysisRefreshJob recover:%v \n%s", err, strace)
				analysisRefreshJob()
			}
		}()
		hourTrigger(func() {
			freedom.ServiceLocator().Call(func(service *domain.AnalysisService) {
				service.ExecuteRefreshJob()
			})
		}, conf.Get().System.JobAnalysisRefreshHour, "analysisRefreshJob")
	}()
}

func dayTrigger(f func(), hour int, name string) {
	getNextTime := func() time.Time {
		addHour := time.Hour * time.Duration(hour)
		date := time.Now().Format("2006-01-02")
		currentDateTime, _ := time.ParseInLocation("2006-01-02", date, time.Local)
		nextTime := currentDateTime.AddDate(0, 0, 1).Add(addHour) //凌晨
		return nextTime
	}

	nextTime := getNextTime()
	freedom.Logger().Infof("%s next time %v", name, nextTime)
	for {
		time.Sleep(15 * time.Second)
		if time.Now().Unix() < nextTime.Unix() {
			continue
		}
		freedom.Logger().Infof("%s next time %v", name, nextTime)
		f()
		nextTime = getNextTime()
	}
}

func hourTrigger(f func(), hour int, name string) {
	getNextTime := func() time.Time {
		addHour := time.Hour * time.Duration(hour)
		date := time.Now().Format("2006-01-02 15")
		currentDateTime, _ := time.ParseInLocation("2006-01-02 15", date, time.Local)
		nextTime := currentDateTime.Add(addHour) //小时
		return nextTime
	}

	nextTime := getNextTime()
	freedom.Logger().Infof("%s next time %v", name, nextTime)
	for {
		time.Sleep(15 * time.Second)
		if time.Now().Unix() < nextTime.Unix() {
			continue
		}
		freedom.Logger().Infof("%s next time %v", name, nextTime)
		f()
		nextTime = getNextTime()
	}
}
