package repository

import (
	"encoding/base64"
	"time"

	"cdp-service/domain/entity"
	"cdp-service/domain/po"
	"cdp-service/infra"
	"cdp-service/infra/cattle"

	"github.com/8treenet/freedom"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *PersonaRepository {
			return &PersonaRepository{}
		})
	})
}

// PersonaRepository 画像仓库.
type PersonaRepository struct {
	freedom.Repository
	Manager *cattle.Manager
	Common  *infra.CommonRequest
}

func (repo *PersonaRepository) Find(id int) (result *entity.Persona, e error) {
	defer func() {
		if e != nil {
			return
		}
		result.XMLBytes, e = base64.StdEncoding.DecodeString(result.XMLData)
	}()
	if id == 0 {
		e = gorm.ErrRecordNotFound
		return
	}
	result = &entity.Persona{}
	result.ID = id
	repo.InjectBaseEntity(result)
	e = findPersona(repo, &result.Persona)
	return
}

func (repo *PersonaRepository) FindByName(name string) (result *entity.Persona, e error) {
	defer func() {
		if e != nil {
			return
		}
		result.XMLBytes, e = base64.StdEncoding.DecodeString(result.XMLData)
	}()
	if name == "" {
		e = gorm.ErrRecordNotFound
		return
	}

	result = &entity.Persona{}
	result.Name = name
	repo.InjectBaseEntity(result)
	e = findPersona(repo, &result.Persona)
	return
}

func (repo *PersonaRepository) NewPersonaEntity() (result *entity.Persona) {
	result = &entity.Persona{}
	result.Created = time.Now()
	result.Updated = time.Now()
	repo.InjectBaseEntity(result)
	return
}

func (repo *PersonaRepository) SavePersonaEntity(entity *entity.Persona) error {
	entity.XMLData = base64.StdEncoding.EncodeToString(entity.XMLBytes)
	if entity.ID == 0 {
		_, e := createPersona(repo, &entity.Persona)
		return e
	}

	_, e := savePersona(repo, entity)
	return e
}

// GetAllPersona
func (repo *PersonaRepository) GetAllPersona() (result []*entity.Persona, e error) {
	result = make([]*entity.Persona, 0)
	list, e := findPersonaList(repo, po.Persona{})
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		entity := &entity.Persona{Persona: list[i]}
		entity.XMLBytes, e = base64.StdEncoding.DecodeString(entity.XMLData)
		if e != nil {
			return
		}

		result = append(result, entity)
	}
	repo.InjectBaseEntitys(result)
	return
}

// GetPersonasByPage .
func (repo *PersonaRepository) GetPersonasByPage() (result []*entity.Persona, totalPage int, e error) {
	result = make([]*entity.Persona, 0)

	page, pageSize := repo.Common.GetPage()
	pager := NewDescPager("id").SetPage(page, pageSize)
	list, e := findPersonaList(repo, po.Persona{}, pager)
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		entity := &entity.Persona{Persona: list[0]}
		entity.XMLBytes, e = base64.StdEncoding.DecodeString(entity.XMLData)
		if e != nil {
			return
		}

		result = append(result, entity)
	}
	totalPage = pager.TotalPage()
	repo.InjectBaseEntitys(result)
	return
}

// BatchCustomerProfile
func (repo *PersonaRepository) BatchCustomerProfile(list []*po.CustomerProfile) error {
	return repo.db().Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(list, 500).Error
}

// DeleteCustomerProfile
func (repo *PersonaRepository) DeleteCustomerProfile(entity *entity.Persona, users []string) error {
	return repo.db().Where("personaId = ? and userId in (?)", entity.ID, users).Delete(&po.CustomerProfile{}).Error
}

// DeleteCustomerProfileByPersonaId
func (repo *PersonaRepository) DeleteCustomerProfileByPersona(entity *entity.Persona) error {
	return repo.db().Where("personaId = ?", entity.ID).Delete(&po.CustomerProfile{}).Error
}

// db .
func (repo *PersonaRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
