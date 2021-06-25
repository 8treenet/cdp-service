// Package repository generated by 'freedom new-project github.com/8treenet/cdp-service'
package repository

import (
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *SupportRepository {
			return &SupportRepository{sourceCacheKey: "cdp_support_source"}
		})
	})
}

// SupportRepository .
type SupportRepository struct {
	freedom.Repository
	sourceCacheKey string
}

// CreateSouce .
func (repo *SupportRepository) CreateSouce(source string) error {
	_, e := createSource(repo, &po.Source{Source: source, Created: time.Now(), Updated: time.Now()})
	return e
}

func (repo *SupportRepository) FindSource(source string) int {
	if source == "" {
		return 0
	}
	var obj po.Source
	key := repo.sourceCacheKey + ":" + source
	if err := redisJSONGet(repo.Redis(), key, &obj); err != nil && err != redis.Nil {
		return obj.ID
	}

	obj.Source = source
	if findSource(repo, &obj) != nil {
		return 0
	}

	redisJSONSet(repo.Redis(), key, obj, time.Minute*100)
	return obj.ID
}

// GetAllSource .
func (repo *SupportRepository) GetAllSource() ([]*po.Source, error) {
	list, err := findSourceList(repo, po.Source{})
	if err != nil {
		return nil, err
	}

	result := []*po.Source{}
	for i := 0; i < len(list); i++ {
		result = append(result, &list[i])
	}
	return result, nil
}

// db .
func (repo *SupportRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
