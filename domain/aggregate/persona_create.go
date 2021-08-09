package aggregate

import (
	"fmt"
	"time"

	"cdp-service/adapter/repository"
	"cdp-service/domain/entity"
	"cdp-service/infra/cattle"
)

// PersonaCreate
type PersonaCreate struct {
	entity.Persona
	featureRepository     *repository.FeatureRepository
	personaRepository     *repository.PersonaRepository
	dataManagerRepository *repository.DataManagerRepository
}

// Do .
func (cmd *PersonaCreate) Do() (e error) {
	dsls, e := cattle.NewArrayDSL(cmd.XMLBytes)
	if e != nil {
		return
	}

	for _, dsl := range dsls {
		tableName := dsl.FindFromNode().GetContent()
		mainFeatureEntity, err := cmd.featureRepository.GetFeatureEntityByWarehouse(tableName)
		if err != nil {
			e = fmt.Errorf("PersonaCreate 未找到表 %s : %w", tableName, err)
			return
		}
		dsl.SetMetedata(mainFeatureEntity)

		jnodes := dsl.FindJoinFromNodes()
		for _, node := range jnodes {
			joinFeature, err := cmd.featureRepository.GetFeatureEntityByWarehouse(node.GetContent())
			if err != nil {
				e = fmt.Errorf("PersonaCreate 未找到关联的表 %s :%w", node.GetContent(), err)
				return
			}
			dsl.SetMetedata(joinFeature)
		}

		builder, err := cattle.ExplainPersonasAnalysis(dsl, []string{"PersonasTest"}, time.Now().Add(1*time.Minute))
		if err != nil {
			e = fmt.Errorf("PersonaCreate ExplainPersonasAnalysis %w", err)
			return
		}
		temp := map[string]interface{}{}
		if e = cmd.dataManagerRepository.Query(&temp, builder); e != nil {
			return
		}
	}

	e = cmd.personaRepository.SavePersonaEntity(&cmd.Persona)
	return
}
