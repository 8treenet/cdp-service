package aggregate

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/cdp-service/utils"
	"github.com/go-xorm/builder"
)

const (
	AnalysisMultipleOutType = "multipleOut"
	AnalysisSingleOutType   = "singleOut"
)

// AnalysisCreate
type AnalysisCreate struct {
	entity.Analysis
	newError              error
	feature               *entity.Feature
	featureRepository     *repository.FeatureRepository
	analysisRepository    *repository.AnalysisRepository
	dataManagerRepository *repository.DataManagerRepository
}

// Do .
func (cmd *AnalysisCreate) Do() (e error) {
	if cmd.newError != nil {
		e = cmd.newError
		return
	}
	if !utils.InSlice([]string{AnalysisMultipleOutType, AnalysisSingleOutType}, cmd.OutType) {
		e = fmt.Errorf("OutType错误 :%s", cmd.OutType)
		return
	}

	dsl, e := cattle.NewDSL(cmd.XMLBytes)
	if e != nil {
		return
	}

	if denominatorNode := dsl.FindDenominatorNode(); denominatorNode != nil {
		//寻找分母
		dentity, err := cmd.analysisRepository.FindByName(denominatorNode.GetContent())
		if err != nil || dentity.OutType != AnalysisSingleOutType || dentity.DenominatorID != 0 {
			e = fmt.Errorf("分母设置%s错误，请检查规则。", denominatorNode.GetContent())
			return
		}
		cmd.DenominatorID = dentity.ID
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

	//检查是否可以查询
	var b *builder.Builder
	if cmd.OutType == AnalysisSingleOutType {
		b, e = cattle.ExplainMultipleAnalysis(dsl, time.Now().Add(1*time.Minute), time.Now())
	} else {
		b, e = cattle.ExplainSingleAnalysis(dsl, time.Now().Add(1*time.Minute), time.Now())
	}
	if e != nil {
		return
	}

	temp := []map[string]interface{}{}
	if e = cmd.dataManagerRepository.Query(&temp, b); e != nil {
		return
	}

	e = cmd.analysisRepository.SaveAnalysisEntity(&cmd.Analysis)
	return
}
