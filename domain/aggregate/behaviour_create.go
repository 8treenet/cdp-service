package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
)

// BehaviourCreate
type BehaviourCreate struct {
	entity.Feature
	behaviours          []*entity.Behaviour
	BehaviourRepository *repository.BehaviourRepository
	DataRepository      *repository.DataRepository
}

// Do .
func (cmd *BehaviourCreate) Do() (e error) {
	successIds := []int{}
	submit := cmd.DataRepository.NewSubmit(cmd.Warehouse)
	m := cmd.ToColumns()
	for k, v := range m {
		submit.AddMetadata(k, v)
	}

	for _, behaviour := range cmd.behaviours {
		dataMap, err := behaviour.ToColumns()
		if err != nil {
			cmd.Worker().Logger().Error(err)
		}

		submit.AddRow(dataMap)
		successIds = append(successIds, behaviour.ID)
	}
	e = cmd.BehaviourRepository.BehavioursSuccess(successIds)
	if e != nil {
		return
	}

	return cmd.DataRepository.SaveSubmit(submit)
}
