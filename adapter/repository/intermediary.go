package repository

import (
	"time"

	"cdp-service/domain/entity"
	"cdp-service/domain/po"

	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *IntermediaryRepository {
			return &IntermediaryRepository{customerMetaDataCacheKey: "cdp_customer_metaData"}
		})
	})
}

// IntermediaryRepository 客户中介仓库 .
type IntermediaryRepository struct {
	freedom.Repository
	customerMetaDataCacheKey string
}

// GetEntity .
func (repo *IntermediaryRepository) GetEntity() (result *entity.Intermediary, e error) {
	templetes, err := repo.getMetaData()
	if err != nil {
		return nil, err
	}

	result = &entity.Intermediary{Templetes: templetes}
	repo.InjectBaseEntity(result)
	return
}

// AddMetaData .
func (repo *IntermediaryRepository) AddMetaData(data po.CustomerExtensionMetadata) error {
	defer func() {
		if e := repo.Redis().Del(repo.customerMetaDataCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()
	data.Created = time.Now()
	data.Updated = data.Created
	if data.Sort == 0 {
		data.Sort = 1000
	}
	_, err := createCustomerExtensionMetadata(repo, &data)
	return err
}

// GetMetaData .
func (repo *IntermediaryRepository) GetMetaData() ([]*po.CustomerExtensionMetadata, error) {
	templetes, err := repo.getMetaData()
	if err != nil {
		return nil, err
	}
	return templetes, nil
}

// GetOneMetaData .
func (repo *IntermediaryRepository) GetOneMetaData(id int) (*po.CustomerExtensionMetadata, error) {
	tmpl := &po.CustomerExtensionMetadata{ID: id}
	return tmpl, findCustomerExtensionMetadata(repo, tmpl)
}

// SaveMetaData .
func (repo *IntermediaryRepository) SaveMetaData(tmpl *po.CustomerExtensionMetadata) error {
	defer func() {
		if e := repo.Redis().Del(repo.customerMetaDataCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()

	_, err := saveCustomerExtensionMetadata(repo, tmpl)
	return err
}

func (repo *IntermediaryRepository) getMetaData() (result []*po.CustomerExtensionMetadata, err error) {
	err = redisJSONGet(repo.Redis(), repo.customerMetaDataCacheKey, &result)
	if err == nil || err != redis.Nil {
		return
	}

	list, err := findCustomerExtensionMetadataList(repo, po.CustomerExtensionMetadata{}, NewDescPager("sort", "id"))
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		result = append(result, &list[i])
	}
	err = redisJSONSet(repo.Redis(), repo.customerMetaDataCacheKey, result)
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
