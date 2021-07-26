package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindFactory(func() *PersonaFactory {
			return &PersonaFactory{}
		})
	})
}

// PersonaFactory
type PersonaFactory struct {
	Worker                freedom.Worker
	FeatureRepository     *repository.FeatureRepository
	PersonaRepository     *repository.PersonaRepository
	DataManagerRepository *repository.DataManagerRepository
}

// CreateAnalysisCMD
func (factory *PersonaFactory) CreatePersonaCMD(name, title string, dateRange int, xmlData []byte) *PersonaCreate {
	result := &PersonaCreate{
		featureRepository:     factory.FeatureRepository,
		personaRepository:     factory.PersonaRepository,
		dataManagerRepository: factory.DataManagerRepository,
	}

	result.Persona = *(factory.PersonaRepository.NewPersonaEntity())
	result.XMLBytes = xmlData
	result.Name = name
	result.Title = title
	result.DateRange = dateRange
	return result
}

// JobPersonaCMD
func (factory *PersonaFactory) JobPersonaCmds() (result []*PersonaJob, e error) {
	all, e := factory.PersonaRepository.GetAllPersona()
	if e != nil {
		return nil, e
	}
	for i := 0; i < len(all); i++ {
		result = append(result, &PersonaJob{
			Persona:               *all[i],
			featureRepository:     factory.FeatureRepository,
			personaRepository:     factory.PersonaRepository,
			dataManagerRepository: factory.DataManagerRepository,
		})
	}
	return
}
