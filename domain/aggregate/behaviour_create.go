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
		behaviour.SyncSuccess()
	}

	defer func() {
		if e != nil {
			return
		}
		err := cmd.DataRepository.SaveSubmit(wholeFlowSubmit)
		if err != nil {
			cmd.Worker().Logger().Error("全站流量提交失败 error:", err)
		}
	}()

	e = cmd.DataRepository.SaveSubmit(submit)
	if e == nil {
		e = cmd.BehaviourRepository.BehavioursFinish(cmd.behaviours)
		return
	}

	for _, v := range cmd.behaviours {
		v.SyncError()
	}
	e = cmd.BehaviourRepository.BehavioursFinish(cmd.behaviours)
	return
}
