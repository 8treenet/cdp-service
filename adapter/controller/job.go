package controller

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"cdp-service/domain"
	"cdp-service/domain/po"
	"cdp-service/server/conf"
	"cdp-service/utils"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(i freedom.Initiator) {
		i.BindBooting(func(bootManager freedom.BootManager) {
			// behaviourSaveJob()
			// behaviourEnteringHouseJob()
			// truncateJob()
			// analysisDayJob()
			// analysisRefreshJob()
			// personaDayJob()
			// personaRefreshJob()
			Press()
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

func personaDayJob() {
	if conf.Get().System.JobPersonaHour == 0 {
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("personaDayJob recover:%v \n%s", err, strace)
				personaDayJob()
			}
		}()
		dayTrigger(func() {
			freedom.ServiceLocator().Call(func(service *domain.PersonaService) {
				service.ExecuteDayJob()
			})
		}, conf.Get().System.JobPersonaHour, "personaDayJob")
	}()
}

func personaRefreshJob() {
	if conf.Get().System.JobPersonaRefreshHour == 0 {
		return
	}

	go func() {
		time.Sleep(1 * time.Second)
		defer func() {
			if err := recover(); err != nil {
				strace := string(utils.GetStackTrace())
				freedom.Logger().Errorf("personaRefreshJob recover:%v \n%s", err, strace)
				personaRefreshJob()
			}
		}()
		hourTrigger(func() {
			freedom.ServiceLocator().Call(func(service *domain.PersonaService) {
				service.ExecuteRefreshJob()
			})
		}, conf.Get().System.JobPersonaRefreshHour, "personaRefreshJob")
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

func Press() {
	go func() {
		time.Sleep(5 * time.Second)
		jdata, _ := json.Marshal(map[string]interface{}{"haha": 20, "hahaname": "btree"})
		freedom.ServiceLocator().Call(func(service *domain.BehaviourService) {
			fmt.Println("kaishi", time.Now().Unix())
			for i := 0; i < 3; i++ {
				var list []*po.Behaviour
				for j := 0; j < 1000; j++ {
					list = append(list, &po.Behaviour{
						UserKey:    fmt.Sprintf("key%d-%d", i, j),
						UserIPAddr: "192.168.1.1",
						FeatureID:  rand.Intn(200),
						CreateTime: time.Now(),
						Data:       jdata,
						Processed:  rand.Intn(3),
						SourceID:   0,
						Created:    time.Now(),
					})
				}
				service.BehaviourRepository.BatchSave(list)
			}
			fmt.Println("jieshu", time.Now().Unix())
		})
	}()
}
