package aggregate

import (
	"encoding/json"
	"fmt"
	"time"

	"cdp-service/adapter/repository"
	"cdp-service/domain/entity"
	"cdp-service/infra/cattle"
	"cdp-service/utils"

	"github.com/go-xorm/builder"
	"gorm.io/gorm"
)

// AnalysisJob 负责分析的处理，并且写入数据库.
type AnalysisJob struct {
	entity.Analysis
	newError              error
	feature               *entity.Feature
	featureRepository     *repository.FeatureRepository
	analysisRepository    *repository.AnalysisRepository
	dataManagerRepository *repository.DataManagerRepository
}

// Do .
func (cmd *AnalysisJob) Do(now time.Time) (e error) {
	if cmd.newError != nil {
		e = cmd.newError
		return
	}

	dsl, e := cattle.NewDSL(cmd.XMLBytes)
	if e != nil {
		return
	}

	//读取元数据来适配类型
	jnodes := dsl.FindJoinFromNodes()
	for _, node := range jnodes {
		joinFeature, err := cmd.featureRepository.GetFeatureEntityByWarehouse(node.GetContent())
		if err != nil {
			e = fmt.Errorf("未找到关联的表 %s :%w", node.GetContent(), err)
			return
		}
		dsl.SetMetedata(joinFeature)
	}
	dsl.SetMetedata(cmd.feature)

	beginTime := cmd.GetBeginTime(now)
	endTime := cmd.GetEndTime(now)
	//解析查询
	var b *builder.Builder
	if cmd.OutType == AnalysisSingleOutType {
		b, e = cattle.ExplainSingleAnalysis(dsl, beginTime, endTime)
	} else {
		b, e = cattle.ExplainMultipleAnalysis(dsl, beginTime, endTime)
	}

	detail := []map[string]interface{}{}
	if e = cmd.dataManagerRepository.Query(&detail, b); e != nil {
		return
	}

	if len(detail) == 0 {
		e = fmt.Errorf("AnalysisTask Query 结果返回0")
		return
	}

	var report *entity.Report
	for {
		report, e = cmd.analysisRepository.GetReportEntity(cmd.ID)
		if e != nil && e != gorm.ErrRecordNotFound {
			e = fmt.Errorf("GetReportEntity %w", e)
			return
		}
		if e == nil {
			break
		}
		report = cmd.analysisRepository.NewReportEntity()
		report.AnalysisID = cmd.ID
		break
	}

	var reportData []byte
	if cmd.OutType == AnalysisSingleOutType {
		reportData, e = json.Marshal(detail[0])
	} else {
		reportData, e = json.Marshal(cmd.conversionList(detail))
	}
	if e != nil {
		return
	}

	report.SetData(reportData)
	report.SetBeginTime(beginTime)
	report.SetEndTime(endTime)
	return cmd.analysisRepository.SaveReportEntity(report)
}

func (cmd *AnalysisJob) conversionList(list []map[string]interface{}) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, vmap := range list {
		newMap := map[string]interface{}{}
		for key, v := range vmap {
			switch v := v.(type) {
			case time.Time:
				newMap[key] = utils.DateTimeFormat(v)
				continue
			}
			newMap[key] = v
		}
		result = append(result, newMap)
	}
	return result
}
