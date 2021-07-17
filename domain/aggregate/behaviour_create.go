package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/infra/cattle"
)

const (
	wholeFlowTableName         = "whole_flow"
	wholeFlowMetadataFeatureId = "featureId"
)

// BehaviourCreate
type BehaviourCreate struct {
	entity.Feature
	WholeFlow           *entity.Feature
	behaviours          []*entity.Behaviour
	BehaviourRepository *repository.BehaviourRepository
	DataRepository      *repository.DataRepository
}

// Do .
func (cmd *BehaviourCreate) Do() (e error) {
	successIds := []int{}
	submit := cmd.DataRepository.NewSubmit(cmd.Warehouse)               //行为提交
	wholeFlowSubmit := cmd.DataRepository.NewSubmit(wholeFlowTableName) //全站流量提交
	wholeFlowSubmit.AddMetadata(wholeFlowMetadataFeatureId, cattle.ColumnTypeUInt16)

	for k, v := range cmd.ToColumns() {
		submit.AddMetadata(k, v)
	}

	for _, behaviour := range cmd.behaviours {
		dataMap, err := behaviour.ToColumns()
		if err != nil {
			cmd.Worker().Logger().Error(err)
		}
		submit.AddRow(behaviour.ID, dataMap)

		dataMap[wholeFlowMetadataFeatureId] = cmd.ID //只需要知道哪个行为的全站流量
		wholeFlowSubmit.AddRow(behaviour.ID, dataMap)
		successIds = append(successIds, behaviour.ID)
	}

	defer func() {
		if e != nil {
			cmd.BehaviourRepository.BehavioursError(successIds) //错误后重置
			return
		}
		cmd.BehaviourRepository.BehavioursSuccess(successIds) //成功后重置
		err := cmd.DataRepository.SaveSubmit(wholeFlowSubmit)
		if err != nil {
			cmd.Worker().Logger().Error("全站流量提交失败 error:", err)
		}
	}()

	e = cmd.DataRepository.SaveSubmit(submit)
	return
}
