package aggregate

import (
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom/infra/transaction"
)

// CustomerCreateCmd
type CustomerCreateCmd struct {
	entity.Intermediary
	CustomerRepo *repository.CustomerRepository //客户仓库
	SignRepo     *repository.SignRepository     //识别仓库
	TX           transaction.Transaction        //依赖倒置事务组件
}

// Do .
func (cmd *CustomerCreateCmd) Do(customerDto vo.CustomerDTO) (e error) {
	customer := cmd.CustomerRepo.CreateCustomer()
	customer.Customer = customerDto.Customer
	customer.Customer.Created = time.Now()
	customer.Customer.Updated = time.Now()
	customer.SetExtension(customerDto.Extension)

	if e = cmd.VerifyCustomer(customer, true); e != nil {
		return
	}

	return cmd.TX.Execute(func() error {
		//获取已经存在的客户
		existCustomer, e := cmd.sign(customerDto)
		if e != nil {
			return e
		}
		if existCustomer != nil {
			return nil
		}
		return cmd.CustomerRepo.SaveCustomer(customer)
	})
}

// BatcheDo .
func (cmd *CustomerCreateCmd) BatcheDo(customerDtos []vo.CustomerDTO) (e error) {
	for _, v := range customerDtos {
		if e = cmd.Do(v); e != nil {
			return
		}
	}
	return
}

// getCustomerByKey .
func (cmd *CustomerCreateCmd) getKey(key string) *po.CustomerKey {
	if key == "" {
		return nil
	}
	result, err := cmd.SignRepo.GetKey(key)
	if err == nil {
		return result
	}
	return nil
}

// getCustomerByPhone .
func (cmd *CustomerCreateCmd) getPhone(phone string) *po.CustomerPhone {
	if phone == "" {
		return nil
	}
	result, err := cmd.SignRepo.GetPhone(phone)
	if err == nil {
		return result
	}
	return nil
}

// getWechat .
func (cmd *CustomerCreateCmd) getWechat(unionId string) *po.CustomerWechat {
	if unionId == "" {
		return nil
	}
	result, err := cmd.SignRepo.GetWechat(unionId)
	if err == nil {
		return result
	}
	return nil
}

func (cmd *CustomerCreateCmd) sign(customerDto vo.CustomerDTO) (*entity.Customer, error) {
	key := cmd.getKey(customerDto.Key)
	phone := cmd.getPhone(customerDto.Phone)
	wechat := cmd.getWechat(customerDto.WechatUnionID)

	if key == nil && phone == nil && wechat == nil {
		return nil, nil
	}
	userId := ""
	if key != nil {
		userId = key.UserID
	}
	if phone != nil && userId == "" {
		userId = phone.UserID
	}
	if wechat != nil && userId == "" {
		userId = wechat.UserID
	}

	if customerDto.Key != "" && key != nil {
		key.SetUserID(userId)
		fmt.Println("fuck")
		cmd.SignRepo.SaveKey(key)
	}
	if customerDto.Key != "" && key == nil {
		obj := po.CustomerKey{UserID: userId, Key: customerDto.Key}
		cmd.SignRepo.SaveKey(&obj)
	}

	if customerDto.Phone != "" && phone != nil {
		phone.SetUserID(userId)
		cmd.SignRepo.SavePhone(phone)
	}
	if customerDto.Phone != "" && phone == nil {
		obj := po.CustomerPhone{UserID: userId, Phone: customerDto.Phone}
		cmd.SignRepo.SavePhone(&obj)
	}

	if customerDto.WechatUnionID != "" && wechat != nil {
		wechat.SetUserID(userId)
		cmd.SignRepo.SaveWechat(wechat)
	}
	if customerDto.WechatUnionID != "" && wechat == nil {
		obj := po.CustomerWechat{UserID: userId, UnionID: customerDto.WechatUnionID}
		cmd.SignRepo.SaveWechat(&obj)
	}

	return cmd.CustomerRepo.GetCustomer(userId)
}
