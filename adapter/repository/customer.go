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
	SignRepository
	Common *infra.CommonRequest
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

func (repo *CustomerRepository) insertCustomer(customer *entity.Customer) error {
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

	if customer.UserKey != "" {
		if err := repo.SaveKey(&po.CustomerKey{UserKey: customer.UserKey, UserID: uuid}); err != nil {
			return err
		}
	}

	if customer.Phone != "" {
		if err := repo.SavePhone(&po.CustomerPhone{Phone: customer.Phone, UserID: uuid}); err != nil {
			return err
		}
	}

	if customer.WechatUnionID != "" {
		if err := repo.SaveWechat(&po.CustomerWechat{UnionID: customer.WechatUnionID, UserID: uuid}); err != nil {
			return err
		}
	}
	return nil
}

// SaveCustomer .
func (repo *CustomerRepository) SaveCustomer(customer *entity.Customer) error {
	if customer.UserID == "" {
		return repo.insertCustomer(customer)
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
	if e = repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerExtension{}).Error; e != nil {
		return
	}
	if entity.Phone != "" {
		repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerPhone{})
	}

	if entity.WechatUnionID != "" {
		repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerWechat{})
	}

	if entity.UserKey != "" {
		repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerKey{})
	}

	return repo.db().Where("userId = ?", entity.UserID).Delete(&po.CustomerExtension{}).Error
}

// GetCustomerKey .
func (repo *CustomerRepository) GetCustomerByKey(key string) (result *entity.Customer, e error) {
	obj := po.CustomerKey{UserKey: key}
	if e = findCustomerKey(repo, &obj); e != nil {
		return
	}
	return repo.GetCustomer(obj.UserID)
}

// GetCustomersByKeys .
func (repo *CustomerRepository) GetCustomersByKey(keys []string) (result []*entity.Customer, e error) {
	objs, e := findCustomerKeyListByWhere(repo, "userKey in (?)", []interface{}{keys})
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
func (repo *CustomerRepository) GetCustomersByPhone(phones []string) (result []*entity.Customer, e error) {
	objs, e := findCustomerPhoneListByWhere(repo, "phone in (?)", []interface{}{phones})
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

// GetCustomersByWechat .
func (repo *CustomerRepository) GetCustomersByWechat(unionIds []string) (result []*entity.Customer, e error) {
	objs, e := findCustomerWechatListByWhere(repo, "unionId in (?)", []interface{}{unionIds})
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

	list, e := findCustomerListByWhere(repo, "userId in (?)", []interface{}{userIdList})
	if e != nil {
		return
	}

	extensions, e := findCustomerExtensionListByWhere(repo, "userId in (?)", []interface{}{userIdList})
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

// GetCustomer .
func (repo *CustomerRepository) GetCustomersByPage() (result []*entity.Customer, totalPage int, e error) {
	result = make([]*entity.Customer, 0)

	page, pageSize := repo.Common.GetPage()
	pager := NewDescPager("id").SetPage(page, pageSize)
	list, e := findCustomerList(repo, po.Customer{}, pager)
	if e != nil {
		return
	}

	var extensionConds []string
	for _, customerPO := range list {
		extensionConds = append(extensionConds, customerPO.UserID)
	}

	extensions, e := findCustomerExtensionListByWhere(repo, "userId in (?)", []interface{}{extensionConds})
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
	totalPage = pager.TotalPage()
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
