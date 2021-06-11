package repository

import (
	"encoding/json"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *CustomerRepository {
			return &CustomerRepository{}
		})
	})
}

// CustomerRepository .
type CustomerRepository struct {
	freedom.Repository
	Common *infra.CommonRequest
	//Mongo                    *infra.Mongo
}

// CreateCustomer .
func (repo *CustomerRepository) CreateCustomer() *entity.Customer {
	result := &entity.Customer{}
	repo.InjectBaseEntity(result)
	return result
}

// SaveCustomer .
func (repo *CustomerRepository) SaveCustomer(customer *entity.Customer) error {
	if customer.UserID == 0 {
		if _, err := createCustomer(repo, &customer.Customer); err != nil {
			return err
		}

		expo := &po.CustomerExtension{UserID: customer.Customer.UserID}
		exBytes, err := json.Marshal(customer.GetExtension())
		if err != nil {
			return err
		}
		expo.Data = datatypes.JSON(exBytes)
		if _, err := createCustomerExtension(repo, expo); err != nil {
			return err
		}
		return nil
	}

	if _, e := saveCustomer(repo, customer); e != nil {
		return e
	}

	extMap := customer.GetExtensionChanges()
	if extMap == nil {
		return nil
	}
	exBytes, err := json.Marshal(extMap)
	if err != nil {
		return err
	}
	jsonMap := map[string]interface{}{"data": datatypes.JSON(exBytes)}
	repo.db().Model(&po.CustomerExtension{}).Where("userId = ?", customer.UserID).Updates(jsonMap)
	if _, e := saveCustomerExtension(repo, &customer.Customer); e != nil {
		return e
	}
	return nil
}

// GetCustomer .
func (repo *CustomerRepository) GetCustomer(id int) (result *entity.Customer, e error) {
	result = &entity.Customer{}
	repo.InjectBaseEntity(result)

	pobj := &po.Customer{UserID: id}
	if e = findCustomer(repo, pobj); e != nil {
		return
	}
	peobj := &po.CustomerExtension{UserID: id}
	if e = findCustomerExtension(repo, peobj); e != nil {
		return
	}

	result.Customer = *pobj
	m := map[string]interface{}{}
	if e = json.Unmarshal(peobj.Data, &m); e != nil {
		return
	}
	result.Extension = m
	return
}

// DeleteCustomer .
func (repo *CustomerRepository) DeleteCustomer(entity *entity.Customer) (e error) {
	if e = repo.db().Where(entity.Location()).Delete(&po.Customer{}).Error; e != nil {
		return
	}
	return repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerExtension{}).Error
}

// db .
func (repo *CustomerRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
