package aggregate

import (
	"encoding/json"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
)

// BehaviourCreate
type BehaviourCreate struct {
	entity.Feature
	behaviours          []*entity.Behaviour
	BehaviourRepository *repository.BehaviourRepository
}

// Do .
func (cmd *BehaviourCreate) Do() (e error) {
	successIds := []int{}

	m666 := map[string]interface{}{}
	m666["Feature"] = cmd.View()
	m666["behaviours"] = cmd.behaviours
	jsonData, _ := json.Marshal(m666)
	cmd.Worker().Logger().Info(string(jsonData))

	for _, behaviour := range cmd.behaviours {
		successIds = append(successIds, behaviour.ID)
	}
	e = cmd.BehaviourRepository.BehavioursSuccess(successIds)
	return
}
