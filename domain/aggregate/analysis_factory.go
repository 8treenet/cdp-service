package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindFactory(func() *AnalysisFactory {
			return &AnalysisFactory{}
		})
	})
}

// AnalysisFactory
type AnalysisFactory struct {
	Worker                freedom.Worker
	FeatureRepository     *repository.FeatureRepository
	AnalysisRepository    *repository.AnalysisRepository
	DataManagerRepository *repository.DataManagerRepository
	Termination           *infra.Termination
}

// CreateAnalysisCMD
func (factory *AnalysisFactory) CreateAnalysisCMD(name, title, outType string, featureId, dateRange, dateConservation, denominatorAnalysisId int, xmlData []byte) *AnalysisCreate {
	result := &AnalysisCreate{
		featureRepository:     factory.FeatureRepository,
		analysisRepository:    factory.AnalysisRepository,
		dataManagerRepository: factory.DataManagerRepository,
	}

	result.Analysis = *(factory.AnalysisRepository.NewAnalysisEntity())
	result.XMLBytes = xmlData
	result.FeatureID = featureId
	result.Name = name
	result.Title = title
	result.DateRange = dateRange
	result.OutType = outType
	result.DateConservation = dateConservation
	result.DenominatorID = denominatorAnalysisId

	feature, err := factory.FeatureRepository.GetFeatureEntity(featureId)
	if err != nil {
		result.newError = err
		return result
	}
	result.feature = feature
	return result
}

func (factory *AnalysisFactory) JobCMD(id int) (result *AnalysisJob) {
	result = &AnalysisJob{
		featureRepository:     factory.FeatureRepository,
		analysisRepository:    factory.AnalysisRepository,
		dataManagerRepository: factory.DataManagerRepository,
	}
	analysisEntity, e := factory.AnalysisRepository.Find(id)
	if e != nil {
		result.newError = e
		return
	}
	result.Analysis = *analysisEntity

	feature, err := factory.FeatureRepository.GetFeatureEntity(result.FeatureID)
	if err != nil {
		result.newError = err
		return
	}
	result.feature = feature
	return
}

// BatchJobCMD
func (factory *AnalysisFactory) BatchJobCMD(dateConservation bool) (result []*AnalysisJob, e error) {
	list, e := factory.AnalysisRepository.GetAllAnalysis()
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		if dateConservation && list[i].DateConservation != 0 {
			continue //自然日 && 不是自然日数据
		}
		if !dateConservation && list[i].DateConservation == 0 {
			continue //不是自然日 && 是自然日数据
		}

		task := &AnalysisJob{
			Analysis:              *list[i],
			featureRepository:     factory.FeatureRepository,
			analysisRepository:    factory.AnalysisRepository,
			dataManagerRepository: factory.DataManagerRepository,
		}

		feature, err := factory.FeatureRepository.GetFeatureEntity(task.FeatureID)
		if err != nil {
			task.newError = err
			continue
		}
		task.feature = feature
		result = append(result, task)
	}
	return
}

// Query
func (factory *AnalysisFactory) Query(id int) (result *AnalysisQuery) {
	result = &AnalysisQuery{
		termination: factory.Termination,
	}

	analysisEntity, e := factory.AnalysisRepository.Find(id)
	if e != nil {
		result.newError = e
		return
	}
	result.Analysis = *analysisEntity

	feature, err := factory.FeatureRepository.GetFeatureEntity(result.FeatureID)
	if err != nil {
		result.newError = err
		return
	}
	result.feature = feature

	report, err := factory.AnalysisRepository.GetReportEntity(id)
	if err != nil {
		result.newError = err
		return
	}
	result.report = report

	if result.DenominatorID == 0 {
		return //没有分母直接返回
	}

	dreport, derr := factory.AnalysisRepository.GetReportEntity(result.DenominatorID)
	if derr != nil {
		result.newError = derr
		return
	}
	result.denominatorReport = dreport
	return
}
