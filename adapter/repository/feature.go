// Package repository generated by 'freedom new-project github.com/8treenet/cdp-service'
package repository

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *FeatureRepository {
			return &FeatureRepository{
				featureCacheKey:              "cdp_support_feature:%d",
				featureAllCacheKey:           "cdp_support_all_feature",
				featureCacheKeyfromWarehouse: "cdp_support_feature_warehouse:%s",
			}
		})
	})
}

// FeatureRepository .
type FeatureRepository struct {
	freedom.Repository
	CommonRequest                *infra.CommonRequest
	sourceCacheKey               string
	featureCacheKey              string
	featureCacheKeyfromWarehouse string
	featureAllCacheKey           string
}

// NewFeatureEntity .
func (repo *FeatureRepository) NewFeatureEntity() *entity.Feature {
	pobj := po.BehaviourFeature{Created: time.Now(), Updated: time.Now()}
	entity := &entity.Feature{BehaviourFeature: pobj}

	repo.InjectBaseEntity(entity)
	return entity
}

// SaveFeatureEntity .
func (repo *FeatureRepository) SaveFeatureEntity(entity *entity.Feature) error {
	defer func() {
		if e := repo.Redis().Del(repo.featureAllCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()

	if entity.ID == 0 {
		if _, e := createBehaviourFeature(repo, &entity.BehaviourFeature); e != nil {
			return e
		}
		for _, metadata := range entity.FeatureMetadata {
			metadata.FeatureID = entity.ID
			if _, e := createBehaviourFeatureMetadata(repo, metadata); e != nil {
				return e
			}
		}
		return nil
	}

	defer func() {
		if e := repo.Redis().Del(fmt.Sprintf(repo.featureCacheKey, entity.ID)).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
		if e := repo.Redis().Del(fmt.Sprintf(repo.featureCacheKeyfromWarehouse, entity.Warehouse)).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()

	for _, metadata := range entity.FeatureMetadata {
		if metadata.ID != 0 {
			continue
		}
		metadata.FeatureID = entity.ID

		if _, e := createBehaviourFeatureMetadata(repo, metadata); e != nil {
			return e
		}
	}

	return nil
}

// GetFeatureEntity .
func (repo *FeatureRepository) GetFeatureEntity(featureId int) (result *entity.Feature, err error) {
	key := fmt.Sprintf(repo.featureCacheKey, featureId)
	result = &entity.Feature{}
	repo.InjectBaseEntity(result)
	err = redisJSONGet(repo.Redis(), key, result)

	result.ID = featureId
	if err == nil {
		return
	}

	if err = findBehaviourFeature(repo, &result.BehaviourFeature); err != nil {
		return
	}

	list, err := findBehaviourFeatureMetadataList(repo, po.BehaviourFeatureMetadata{FeatureID: featureId})
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		result.FeatureMetadata = append(result.FeatureMetadata, &list[i])
	}

	redisJSONSet(repo.Redis(), key, result)
	return
}

// GetFeatureEntityByWarehouse .
func (repo *FeatureRepository) GetFeatureEntityByWarehouse(warehouse string) (result *entity.Feature, err error) {
	key := fmt.Sprintf(repo.featureCacheKeyfromWarehouse, warehouse)
	result = &entity.Feature{}
	result.Warehouse = warehouse
	repo.InjectBaseEntity(result)
	err = redisJSONGet(repo.Redis(), key, result)
	if err == nil {
		return
	}

	if err = findBehaviourFeature(repo, &result.BehaviourFeature); err != nil {
		return
	}

	list, err := findBehaviourFeatureMetadataList(repo, po.BehaviourFeatureMetadata{FeatureID: result.ID})
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		result.FeatureMetadata = append(result.FeatureMetadata, &list[i])
	}

	redisJSONSet(repo.Redis(), key, result)
	return
}

// GetAllFeatureEntity .
func (repo *FeatureRepository) GetAllFeatureEntity() (result []*entity.Feature, err error) {
	defer func() {
		if len(result) == 0 {
			return
		}
		repo.InjectBaseEntitys(result)
	}()

	err = redisJSONGet(repo.Redis(), repo.featureAllCacheKey, &result)
	if err == nil || err != redis.Nil {
		return
	}

	entityList, err := findBehaviourFeatureListByWhere(repo, "warehouse not in(?)", []interface{}{"whole_flow"})
	if err != nil {
		return
	}

	featureIds := []int{}
	for i := 0; i < len(entityList); i++ {
		featureIds = append(featureIds, entityList[i].ID)
		result = append(result, &entity.Feature{BehaviourFeature: entityList[i]})
	}

	metadataList, err := findBehaviourFeatureMetadataListByWhere(repo, "featureId in (?)", []interface{}{featureIds})
	if err != nil {
		return
	}
	for i := 0; i < len(metadataList); i++ {
		for _, entity := range result {
			if entity.ID == metadataList[i].FeatureID {
				entity.FeatureMetadata = append(entity.FeatureMetadata, &metadataList[i])
				break
			}
		}
	}
	redisJSONSet(repo.Redis(), repo.featureAllCacheKey, result)
	return
}

// GetFeatureEntitys .
func (repo *FeatureRepository) GetFeatureEntitys() (result []*entity.Feature, total int, err error) {
	defer func() {
		if len(result) == 0 {
			return
		}
		repo.InjectBaseEntitys(result)
	}()

	page, pageSize := repo.CommonRequest.GetPage()

	pager := NewDescPager("id").SetPage(page, pageSize)
	entityList, err := findBehaviourFeatureListByWhere(repo, "warehouse not in(?)", []interface{}{"whole_flow"}, pager)
	if err != nil {
		return
	}

	featureIds := []int{}
	for i := 0; i < len(entityList); i++ {
		featureIds = append(featureIds, entityList[i].ID)
		result = append(result, &entity.Feature{BehaviourFeature: entityList[i]})
	}

	metadataList, err := findBehaviourFeatureMetadataListByWhere(repo, "featureId in (?)", []interface{}{featureIds})
	if err != nil {
		return
	}
	for i := 0; i < len(metadataList); i++ {
		for _, entity := range result {
			if entity.ID == metadataList[i].FeatureID {
				entity.FeatureMetadata = append(entity.FeatureMetadata, &metadataList[i])
				break
			}
		}
	}

	total = pager.TotalPage()
	return
}

// db .
func (repo *FeatureRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
