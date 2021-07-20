package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
)

// AnalysisTask 负责分析的处理，并且写入数据库.
type AnalysisTask struct {
	entity.Analysis
	feature            *entity.Feature
	featureRepository  *repository.FeatureRepository
	analysisRepository *repository.AnalysisRepository
}

// Do .
func (cmd *AnalysisTask) Do() (e error) {
	// dsl, e := cattle.NewDSL(cmd.XMLBytes)
	// if e != nil {
	// 	return
	// }

	return
}

// // Do .
// func (cmd *AnalysisTask) GetSingle() *entity.Report {
// 	dsl, e := cattle.NewDSL(cmd.XMLBytes)
// 	if e != nil {
// 		return
// 	}
// }
