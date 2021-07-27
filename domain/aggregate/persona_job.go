package aggregate

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/cdp-service/utils"
)

// PersonaJob
type PersonaJob struct {
	entity.Persona
	featureRepository     *repository.FeatureRepository
	personaRepository     *repository.PersonaRepository
	dataManagerRepository *repository.DataManagerRepository
}

// Do .
func (cmd *PersonaJob) Do(checkUsers []string, now time.Time) (e error) {
	users, e := cmd.getMatchUsers(checkUsers, now)
	if e != nil {
		return
	}

	var deleteUsers []string
	for _, v := range checkUsers {
		if utils.InSlice(users, v) {
			continue
		}
		deleteUsers = append(deleteUsers, v)
	}

	if len(deleteUsers) > 0 {
		if err := cmd.personaRepository.DeleteCustomerProfile(&cmd.Persona, deleteUsers); err != nil {
			cmd.Worker().Logger().Errorf("DeleteCustomerProfile失败: entity:%v err:%v", *cmd, err)
		}
	}

	var profiles []*po.CustomerProfile
	for _, userId := range users {
		profiles = append(profiles, &po.CustomerProfile{
			UserID:    userId,
			PersonaID: cmd.ID,
			BeginTime: cmd.GetBeginTime(now),
			EndTime:   now,
			Created:   time.Now(),
			Updated:   time.Now(),
		})
	}

	return cmd.personaRepository.BatchCustomerProfile(profiles)
}

func (cmd *PersonaJob) getMatchUsers(checkUsers []string, now time.Time) (users []string, e error) {
	userMap := map[string]struct{}{}
	for _, user := range checkUsers {
		userMap[user] = struct{}{}
	}

	dsls, e := cattle.NewArrayDSL(cmd.XMLBytes)
	if e != nil {
		return
	}

	for _, dsl := range dsls {
		tableName := dsl.FindFromNode().GetContent()
		mainFeatureEntity, err := cmd.featureRepository.GetFeatureEntityByWarehouse(tableName)
		if err != nil {
			e = fmt.Errorf("PersonaJob 未找到表 %s : %w", tableName, err)
			return
		}
		dsl.SetMetedata(mainFeatureEntity)

		jnodes := dsl.FindJoinFromNodes()
		for _, node := range jnodes {
			joinFeature, err := cmd.featureRepository.GetFeatureEntityByWarehouse(node.GetContent())
			if err != nil {
				e = fmt.Errorf("PersonaJob 未找到关联的表 %s :%w", node.GetContent(), err)
				return
			}
			dsl.SetMetedata(joinFeature)
		}

		builder, err := cattle.ExplainPersonasAnalysis(dsl, checkUsers, cmd.GetBeginTime(now))
		if err != nil {
			e = fmt.Errorf("PersonaJob ExplainPersonasAnalysis %w", err)
			return
		}

		matchUsers := []string{}
		if e = cmd.dataManagerRepository.Query(&matchUsers, builder); e != nil {
			return
		}

		for k := range userMap {
			if utils.InSlice(matchUsers, k) {
				continue
			}
			delete(userMap, k)
		}
	}

	for key := range userMap {
		users = append(users, key)
	}
	return
}
