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
func (factory *AnalysisFactory) CreateAnalysisCMD(name, title, outType string, featureId, dateRange int, xmlData []byte) *AnalysisCreate {
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

	feature, err := factory.FeatureRepository.GetFeatureEntity(featureId)
	if err != nil {
		result.newError = err
		return result
	}
	result.feature = feature
	return result
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
