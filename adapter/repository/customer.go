package repository

import (
	"encoding/json"
	"time"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/cdp-service/utils"
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
	Common      *infra.CommonRequest
	Termination *infra.Termination
	//Mongo                    *infra.Mongo
}

// GetCustomer .
func (repo *CustomerRepository) GetCustomer(id string) (result *entity.Customer, e error) {
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

// CreateCustomer .
func (repo *CustomerRepository) CreateCustomer() *entity.Customer {
	result := &entity.Customer{}
	repo.InjectBaseEntity(result)
	return result
}

// SaveCustomer .
func (repo *CustomerRepository) SaveCustomer(customer *entity.Customer) error {
	if customer.UserID == "" {
		uuid, err := utils.GenerateUUID()
		if err != nil {
			return err
		}
		customer.UserID = uuid
		if _, err := createCustomer(repo, &customer.Customer); err != nil {
			return err
		}

		expo := &po.CustomerExtension{UserID: customer.Customer.UserID, Created: time.Now(), Updated: time.Now()}
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

	extMap := customer.GetExtension()
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

// DeleteCustomer .
func (repo *CustomerRepository) DeleteCustomer(entity *entity.Customer) (e error) {
	if e = repo.db().Where(entity.Location()).Delete(&po.Customer{}).Error; e != nil {
		return
	}
	return repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerExtension{}).Error
}

// GetCustomerKey .
func (repo *CustomerRepository) GetCustomerByKey(key string) (result *entity.Customer, e error) {
	obj := po.CustomerKey{Key: key}
	if e = findCustomerKey(repo, &obj); e != nil {
		return
	}
	return repo.GetCustomer(obj.UserID)
}

// GetCustomersByKeys .
func (repo *CustomerRepository) GetCustomersByKeys(keys []string) (result []*entity.Customer, e error) {
	var keysCond []interface{}
	for _, cond := range keys {
		keysCond = append(keysCond, cond)
	}

	objs, e := findCustomerKeyListByWhere(repo, "key in ?", keysCond)
	if e != nil || len(objs) == 0 {
		return
	}

	var userIdList []string
	for _, v := range objs {
		userIdList = append(userIdList, v.UserID)
	}
	return repo.GetCustomers(userIdList)
}

// GetCustomerByPhone .
func (repo *CustomerRepository) GetCustomerByPhone(phone string) (result *entity.Customer, e error) {
	obj := po.CustomerPhone{Phone: phone}
	if e = findCustomerPhone(repo, &obj); e != nil {
		return
	}
	return repo.GetCustomer(obj.UserID)
}

// GetCustomersByPhones .
func (repo *CustomerRepository) GetCustomersByPhones(phones []string) (result []*entity.Customer, e error) {
	var phonesCond []interface{}
	for _, cond := range phones {
		phonesCond = append(phonesCond, cond)
	}

	objs, e := findCustomerPhoneListByWhere(repo, "phone in ?", phonesCond)
	if e != nil || len(objs) == 0 {
		return
	}

	var userIdList []string
	for _, v := range objs {
		userIdList = append(userIdList, v.UserID)
	}
	return repo.GetCustomers(userIdList)
}

// GetCustomerByWechat .
func (repo *CustomerRepository) GetCustomerByWechat(unionId string) (result *entity.Customer, e error) {
	obj := po.CustomerWechat{UnionID: unionId}
	if e = findCustomerWechat(repo, &obj); e != nil {
		return
	}
	return repo.GetCustomer(obj.UserID)
}

// GetCustomersByWechats .
func (repo *CustomerRepository) GetCustomersByWechats(unionIds []string) (result []*entity.Customer, e error) {
	var unionIdsCond []interface{}
	for _, cond := range unionIds {
		unionIdsCond = append(unionIdsCond, cond)
	}

	objs, e := findCustomerWechatListByWhere(repo, "unionId in ?", unionIdsCond)
	if e != nil || len(objs) == 0 {
		return
	}

	var userIdList []string
	for _, v := range objs {
		userIdList = append(userIdList, v.UserID)
	}
	return repo.GetCustomers(userIdList)
}

// GetCustomer .
func (repo *CustomerRepository) GetCustomers(userIdList []string) (result []*entity.Customer, e error) {
	result = make([]*entity.Customer, 0)
	var primarys []interface{}
	for _, primary := range result {
		primarys = append(primarys, primary)
	}

	list, e := findCustomerListByPrimarys(repo, primarys...)
	if e != nil {
		return
	}

	extensions, e := findCustomerExtensionListByWhere(repo, "userId in ?", primarys)
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		centity := &entity.Customer{Customer: list[i], Extension: make(map[string]interface{})}

		for j := 0; j < len(extensions); j++ {
			if extensions[j].UserID != centity.UserID {
				continue
			}
			if e := json.Unmarshal(extensions[j].Data, &centity.Extension); e != nil {
				repo.Worker().Logger().Error("CustomerRepository.GetCustomers ", e)
			}
			break
		}
		result = append(result, centity)
	}

	repo.InjectBaseEntitys(result)
	return
}

// db .
func (repo *CustomerRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
