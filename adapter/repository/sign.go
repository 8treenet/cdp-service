package repository

import (
	"time"

	"cdp-service/domain/po"
	"cdp-service/infra"

	"github.com/8treenet/freedom"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *SignRepository {
			return &SignRepository{}
		})
	})
}

// SignRepository 用户唯一相关.
type SignRepository struct {
	freedom.Repository
	Termination *infra.Termination
}

// GetCustomer .
func (repo *SignRepository) GetWechat(unionId string) (result *po.CustomerWechat, e error) {
	result = &po.CustomerWechat{UnionID: unionId}
	if e = findCustomerWechat(repo, result); e != nil {
		return
	}
	return
}

// GetPhone .
func (repo *SignRepository) GetPhone(phone string) (result *po.CustomerPhone, e error) {
	result = &po.CustomerPhone{Phone: phone}
	if e = findCustomerPhone(repo, result); e != nil {
		return
	}
	return
}

// GetKey .
func (repo *SignRepository) GetKey(key string) (result *po.CustomerKey, e error) {
	result = &po.CustomerKey{UserKey: key}
	if e = findCustomerKey(repo, result); e != nil {
		return
	}
	return
}

// SaveWechat .
func (repo *SignRepository) SaveWechat(obj *po.CustomerWechat) error {
	if obj.ID != 0 {
		_, err := saveCustomerWechat(repo, obj)
		return err
	}

	obj.Created = time.Now()
	obj.Updated = time.Now()
	_, e := createCustomerWechat(repo, obj)
	return e
}

// SavePhone .
func (repo *SignRepository) SavePhone(obj *po.CustomerPhone) error {
	if obj.ID != 0 {
		_, err := saveCustomerPhone(repo, obj)
		return err
	}

	obj.Created = time.Now()
	obj.Updated = time.Now()
	_, e := createCustomerPhone(repo, obj)
	return e
}

// SaveKey .
func (repo *SignRepository) SaveKey(obj *po.CustomerKey) error {
	if obj.ID != 0 {
		_, err := saveCustomerKey(repo, obj)
		return err
	}

	obj.Created = time.Now()
	obj.Updated = time.Now()
	_, e := createCustomerKey(repo, obj)
	return e
}

// db .
func (repo *SignRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
