package domain

import (
	"sort"

	"github.com/8treenet/cdp-service/domain/aggregate"
	"github.com/8treenet/cdp-service/domain/vo"
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
	cmds, err := service.AnalysisFactory.BatchJobCMD(true)
	if err != nil {
		service.Worker.Logger().Errorf("DayTask.CreateBatchTaskCMD err :%s", err.Error())
		return
	}
	for _, cmd := range cmds {
		err := cmd.Do()
		if err == nil {
			continue
		}
		service.Worker.Logger().Errorf("DayTask.Cmd.Do err :%s", err.Error())
	}
}

// ExecuteRefreshJob
func (service *AnalysisService) ExecuteRefreshJob() {
	cmds, err := service.AnalysisFactory.BatchJobCMD(false)
	if err != nil {
		service.Worker.Logger().Errorf("DayTask.CreateBatchTaskCMD err :%s", err.Error())
		return
	}

	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].DenominatorID < cmds[j].DenominatorID //分母在前
	})

	for _, cmd := range cmds {
		err := cmd.Do()
		if err == nil {
			continue
		}
		service.Worker.Logger().Errorf("DayTask.CMD.Do err :%s", err.Error())
	}
}

// CreateAnalysis
func (service *AnalysisService) CreateAnalysis(req vo.ReqCreateAnalysis) error {
	cmd := service.AnalysisFactory.CreateAnalysisCMD(req.Name, req.Title, req.OutType, req.FeatureId, req.DateRange, req.XmlData)
	return cmd.Do()
}

// QueryAnalysis
func (service *AnalysisService) QueryAnalysis(id int) error {
	query := service.AnalysisFactory.Query(id)
	return query.Do()
}
