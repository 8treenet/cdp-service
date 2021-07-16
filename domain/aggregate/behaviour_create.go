package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
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
	submit := cmd.DataRepository.NewSubmit(cmd.Warehouse)         //行为提交
	wholeFlowSubmit := cmd.DataRepository.NewSubmit("whole_flow") //全站流量提交
	wholeFlowSubmit.AddMetadata("feature", "String")

	for k, v := range cmd.ToColumns() {
		submit.AddMetadata(k, v)
	}
	for k, v := range cmd.WholeFlow.ToColumns() {
		wholeFlowSubmit.AddMetadata(k, v)
	}

	for _, behaviour := range cmd.behaviours {
		dataMap, err := behaviour.ToColumns()
		if err != nil {
			cmd.Worker().Logger().Error(err)
		}
		submit.AddRow(dataMap)

		dataMap["feature"] = cmd.Warehouse
		wholeFlowSubmit.AddRow(dataMap)
		successIds = append(successIds, behaviour.ID)
	}

	e = cmd.BehaviourRepository.BehavioursSuccess(successIds)
	if e != nil {
		return
	}

	defer func() {
		err := cmd.DataRepository.SaveSubmit(wholeFlowSubmit)
		if err != nil {
			cmd.Worker().Logger().Error("全站流量提交失败 error:", err)
		}
	}()
	return cmd.DataRepository.SaveSubmit(submit)
}
