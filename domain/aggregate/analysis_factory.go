package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
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
	dataManagerRepository *repository.DataManagerRepository
}

// CreateAnalysisCMD
func (factory *AnalysisFactory) CreateAnalysisCMD(name, title, outType string, featureId, dateRange int, xmlData []byte) *AnalysisCreate {
	result := &AnalysisCreate{
		featureRepository:     factory.FeatureRepository,
		analysisRepository:    factory.AnalysisRepository,
		dataManagerRepository: factory.dataManagerRepository,
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
