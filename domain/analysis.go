package domain

import (
	"encoding/json"
	"time"

	"sort"

	"github.com/8treenet/cdp-service/domain/aggregate"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *AnalysisService {
			return &AnalysisService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *AnalysisService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// AnalysisService .
type AnalysisService struct {
	Worker          freedom.Worker
	AnalysisFactory *aggregate.AnalysisFactory
}

// ExecuteDayJob
func (service *AnalysisService) ExecuteDayJob() {
	now := time.Now()
	cmds, err := service.AnalysisFactory.BatchJobCMD(true)
	if err != nil {
		service.Worker.Logger().Errorf("ExecuteDayJob err :%s", err.Error())
		return
	}
	for _, cmd := range cmds {
		err := cmd.Do(now)
		if err == nil {
			continue
		}
		service.Worker.Logger().Errorf("DayTask.Cmd.Do err :%s", err.Error())
	}
}

// ExecuteRefreshJob
func (service *AnalysisService) ExecuteRefreshJob() {
	now := time.Now()
	cmds, err := service.AnalysisFactory.BatchJobCMD(false)
	if err != nil {
		service.Worker.Logger().Errorf("ExecuteRefreshJob err :%s", err.Error())
		return
	}

	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].DenominatorID < cmds[j].DenominatorID //分母在前
	})

	for _, cmd := range cmds {
		err := cmd.Do(now)
		if err == nil {
			continue
		}
		service.Worker.Logger().Errorf("DayTask.CMD.Do err :%s", err.Error())
	}
}

// CreateAnalysis
func (service *AnalysisService) CreateAnalysis(req vo.ReqCreateAnalysis) error {
	cmd := service.AnalysisFactory.CreateAnalysisCMD(req.Name, req.Title, req.OutType, req.FeatureId, req.DateRange, req.DateConservation, req.DenominatorAnalysisId, req.XmlData)
	err := cmd.Do()
	if err != nil {
		return err
	}

	utils.Async("CreateAnalysis.Job", service.Worker, func() {
		job := service.AnalysisFactory.JobCMD(cmd.ID)
		e := job.Do(time.Now())
		if e != nil {
			service.Worker.Logger().Error("CreateAnalysis.Job id:%d error:%v", cmd.ID, e)
		}
	})
	return nil
}

// QueryAnalysis
func (service *AnalysisService) QueryAnalysis(id int) (json.Marshaler, error) {
	query := service.AnalysisFactory.Query(id)
	if e := query.Do(); e != nil {
		return nil, e
	}
	return query, nil
}
