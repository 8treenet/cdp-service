package aggregate

import (
	"encoding/json"
	"net"
	"time"

	"cdp-service/adapter/repository"
	"cdp-service/domain/entity"
	"cdp-service/domain/po"
	"cdp-service/domain/vo"
	"cdp-service/infra/geo"
	"cdp-service/utils"

	"github.com/8treenet/freedom/infra/transaction"
)

// CustomerCreateCmd
type CustomerCreateCmd struct {
	entity.Intermediary
	CustomerRepo        *repository.CustomerRepository //客户仓库
	SignRepo            *repository.SignRepository     //识别仓库
	TX                  transaction.Transaction        //依赖倒置事务组件
	GEO                 *geo.GEO                       //geo
	UserRegisterEntity  *entity.Feature                //注册行为实体
	SupportRepository   *repository.SupportRepository
	BehaviourRepository *repository.BehaviourRepository
}

// Do .
func (cmd *CustomerCreateCmd) Do(customerDto vo.CustomerDTO, inBatche ...bool) (e error) {
	customer := cmd.CustomerRepo.CreateCustomer()
	customer.Customer = customerDto.Customer
	if customerDto.BirthdaySubstitute != "" {
		birthday, err := time.ParseInLocation("2006-01-02", customerDto.BirthdaySubstitute, time.Local)
		if err == nil {
			customer.Customer.Birthday = &birthday
		}
	}

	customer.Customer.Created = time.Now()
	customer.Customer.Updated = time.Now()
	for customerDto.RegisterDateTime != "" {
		RegisterTime, err := time.ParseInLocation("2006-01-02 15:04:05", customerDto.RegisterDateTime, time.Local)
		if err != nil {
			cmd.Worker().Logger().Error("CustomerCreateCmd :%v,err:%v", customerDto, err)
			break
		}
		customer.Customer.Created = RegisterTime
		break
	}

	customer.SetExtension(customerDto.Extension)
	for len(inBatche) == 0 {
		customer.SourceID = cmd.SupportRepository.FindSourceID(customerDto.Source)
		if customerDto.IP == "" || net.ParseIP(customerDto.IP) == nil || customerDto.City != "" || customerDto.Region != "" {
			break
		}
		geoInfo, err := cmd.GEO.ParseIP(customerDto.IP)
		if err != nil && geoInfo != nil {
			break
		}
		customer.City = geoInfo.City
		customer.Region = geoInfo.Region
		break
	}

	if e = cmd.VerifyCustomer(customer, true); e != nil {
		return
	}

	exist := true
	err := cmd.TX.Execute(func() error {
		//获取已经存在的客户
		existCustomer, e := cmd.sign(customerDto)
		if e != nil {
			return e
		}
		if existCustomer != nil {
			return nil
		}
		exist = false

		return cmd.CustomerRepo.SaveCustomer(customer)
	})
	if err != nil || exist {
		return nil
	}

	b, err := cmd.createRegisterBehaviour(customer, customerDto.IP)
	if err == nil {
		cmd.BehaviourRepository.AddQueue([]*po.Behaviour{b})
	}
	return nil
}

// BatcheDo .
func (cmd *CustomerCreateCmd) BatcheDo(customerDtos []vo.CustomerDTO) (e error) {
	sourceMap := map[string]int{}
	var ipAddrs []string
	allSource, e := cmd.SupportRepository.GetAllSource()
	if e != nil {
		return
	}
	for _, sourcePo := range allSource {
		sourceMap[sourcePo.Source] = sourcePo.ID
	}

	for i := 0; i < len(customerDtos); i++ {
		if customerDtos[i].IP == "" || utils.InSlice(ipAddrs, customerDtos[i].IP) {
			continue
		}
		ipAddrs = append(ipAddrs, customerDtos[i].IP)
	}

	addrsMap, _ := cmd.GEO.ParseBatchIP(ipAddrs)
	for _, v := range customerDtos {
		if info, ok := addrsMap[v.IP]; ok {
			v.City = info.City
			v.Region = info.Region
		}
		v.SourceID = sourceMap[v.Source]
		if e = cmd.Do(v, true); e != nil {
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
	key := cmd.getKey(customerDto.UserKey)
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

	if customerDto.UserKey != "" && key != nil {
		key.SetUserID(userId)
		cmd.SignRepo.SaveKey(key)
	}
	if customerDto.UserKey != "" && key == nil {
		obj := po.CustomerKey{UserID: userId, UserKey: customerDto.UserKey}
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

// createRegisterBehaviour .
func (cmd *CustomerCreateCmd) createRegisterBehaviour(customer *entity.Customer, ip string) (*po.Behaviour, error) {
	data := map[string]interface{}{}
	for k, v := range customer.Extension {
		data[k] = v
	}
	data["userId"] = customer.UserID
	data["name"] = customer.Name
	data["email"] = customer.Email
	data["phone"] = customer.Phone
	data["gender"] = customer.Gender
	data["birthday"] = customer.Birthday.Format("2006-01-02")

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	result := &po.Behaviour{
		ID:            0,
		WechatUnionID: customer.WechatUnionID,
		UserKey:       customer.UserKey,
		UserPhone:     customer.Phone,
		TempUserID:    "",
		UserIPAddr:    ip,
		FeatureID:     cmd.UserRegisterEntity.ID,
		CreateTime:    customer.Created,
		Data:          jsonData,
		Processed:     0,
		SourceID:      customer.SourceID,
		Created:       time.Now(),
	}
	return result, nil
}
