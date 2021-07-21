package aggregate

import (
	"encoding/json"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/infra"
)

// AnalysisQuery 负责分析的查询.
type AnalysisQuery struct {
	entity.Analysis
	newError          error
	feature           *entity.Feature
	report            *entity.Report
	denominatorReport *entity.Report
	termination       *infra.Termination
	result            interface{}
}

// Do .
func (cmd *AnalysisQuery) Do() (e error) {
	if cmd.report == nil {
		return cmd.termination.AnalysisIng()
	}
	if cmd.DenominatorID > 0 && cmd.denominatorReport == nil {
		return cmd.termination.AnalysisIng() //正在分析中
	}

	if cmd.DenominatorID == 0 {
		return json.Unmarshal(cmd.report.Data, cmd.result) //原始数据
	}

	if cmd.OutType == AnalysisSingleOutType {
		cmd.result, e = cmd.denominatorReport.RatioSingle(cmd.report)
		return
	}
	cmd.result, e = cmd.denominatorReport.RatioMultiple(cmd.report)
	return
}

// MarshalJSON .
func (cmd *AnalysisQuery) MarshalJSON() ([]byte, error) {
	return json.Marshal(cmd.result)
}
