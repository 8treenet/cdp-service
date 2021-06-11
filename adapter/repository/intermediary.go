package repository

import (
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *IntermediaryRepository {
			return &IntermediaryRepository{customerTemplateCacheKey: "cdp_customer_template"}
		})
	})
}

// IntermediaryRepository 客户中介仓库 .
type IntermediaryRepository struct {
	freedom.Repository
	customerTemplateCacheKey string
}

// GetEntity .
func (repo *IntermediaryRepository) GetEntity() (result *entity.Intermediary, e error) {
	templetes, err := repo.getTempletes()
	if err != nil {
		return nil, err
	}

	result = &entity.Intermediary{Templetes: templetes}
	repo.InjectBaseEntity(result)
	return
}

// AddTemplete .
func (repo *IntermediaryRepository) AddTemplete(name, kind, dict, reg string, required int) error {
	defer func() {
		if e := repo.Redis().Del(repo.customerTemplateCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()
	pobject := &po.CustomerExtensionTemplate{
		Name:     name,
		Kind:     kind,
		Dict:     dict,
		Reg:      reg,
		Required: required,
		Sort:     1000,
		Created:  time.Now(),
		Updated:  time.Now(),
	}
	_, err := createCustomerExtensionTemplate(repo, pobject)
	return err
}

// GetTempletes .
func (repo *IntermediaryRepository) GetTempletes() ([]*po.CustomerExtensionTemplate, error) {
	templetes, err := repo.getTempletes()
	if err != nil {
		return nil, err
	}
	return templetes, nil
}

// GetTemplete .
func (repo *IntermediaryRepository) GetTemplete(id int) (*po.CustomerExtensionTemplate, error) {
	tmpl := &po.CustomerExtensionTemplate{ID: id}
	return tmpl, findCustomerExtensionTemplate(repo, tmpl)
}

// SaveTemplete .
func (repo *IntermediaryRepository) SaveTemplete(tmpl *po.CustomerExtensionTemplate) error {
	defer func() {
		if e := repo.Redis().Del(repo.customerTemplateCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()

	_, err := saveCustomerExtensionTemplate(repo, tmpl)
	return err
}

func (repo *IntermediaryRepository) getTempletes() (result []*po.CustomerExtensionTemplate, err error) {
	err = redisJSONGet(repo.Redis(), repo.customerTemplateCacheKey, &result)
	if err == nil || err != redis.Nil {
		return
	}

	list, err := findCustomerExtensionTemplateList(repo, po.CustomerExtensionTemplate{}, NewDescPager("sort", "id"))
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		result = append(result, &list[i])
	}
	err = redisJSONSet(repo.Redis(), repo.customerTemplateCacheKey, result)
	return
}

// db .
func (repo *IntermediaryRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
