package aggregate

import (
	"cdp-service/adapter/repository"
	"cdp-service/domain/entity"
	"cdp-service/server/conf"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindFactory(func() *BehaviourFactory {
			return &BehaviourFactory{}
		})
	})
}

// BehaviourFactory 行为工厂
type BehaviourFactory struct {
	Worker              freedom.Worker
	BehaviourRepository *repository.BehaviourRepository
	FeatureRepository   *repository.FeatureRepository
	CustomerRepo        *repository.CustomerRepository
	DataRepository      *repository.DataManagerRepository
}

// CreateBehaviourCmd
func (factory *BehaviourFactory) CreateBehaviourCmds() (cmds []*BehaviourCreate, e error) {
	list, e := factory.FeatureRepository.GetAllFeatureEntity()
	if e != nil {
		return
	}

	for _, featureEntity := range list {
		if featureEntity.SkipWarehouse != 0 {
			continue
		}
		behaviours, e := factory.BehaviourRepository.FetchBehaviours(featureEntity.ID, 0, conf.Get().System.JobEnteringHouseMaxCount)
		if len(behaviours) == 0 {
			continue
		}
		if e != nil {
			factory.Worker.Logger().Error(e)
			continue
		}
		if err := factory.setGEO(behaviours); err != nil {
			factory.Worker.Logger().Error(e)
		}

		factory.setUserId(behaviours)
		cmds = append(cmds, &BehaviourCreate{
			DataRepository:      factory.DataRepository,
			BehaviourRepository: factory.BehaviourRepository,
			Feature:             *featureEntity,
			behaviours:          behaviours,
		})
	}
	return
}

func (factory *BehaviourFactory) setGEO(behaviours []*entity.Behaviour) error {
	var ips []string
	for i := 0; i < len(behaviours); i++ {
		if behaviours[i].UserIPAddr == "" {
			continue
		}
		ips = append(ips, behaviours[i].UserIPAddr)
	}

	ipMap, err := factory.BehaviourRepository.GetIP(ips)
	if err != nil {
		return err
	}

	for i := 0; i < len(behaviours); i++ {
		geo, ok := ipMap[behaviours[i].UserIPAddr]
		if !ok {
			continue
		}
		behaviours[i].Region = geo.Region
		behaviours[i].City = geo.City
	}
	return nil
}

func (factory *BehaviourFactory) setUserId(behaviours []*entity.Behaviour) {
	var (
		wechatUnionIDs []string
		userKeys       []string
		userPhones     []string
	)
	wechatUnionMap := map[string]string{}
	userKeyMap := map[string]string{}
	userPhoneMap := map[string]string{}

	for _, v := range behaviours {
		if v.UserKey != "" {
			userKeys = append(userKeys, v.UserKey)
		}
		if v.WechatUnionID != "" {
			wechatUnionIDs = append(wechatUnionIDs, v.WechatUnionID)
		}
		if v.UserPhone != "" {
			userPhones = append(userPhones, v.UserPhone)
		}
	}

	for len(wechatUnionIDs) > 0 {
		list, err := factory.CustomerRepo.GetCustomersByWechat(wechatUnionIDs)
		if err != nil {
			factory.Worker.Logger().Error(err)
		}
		for i := 0; i < len(list); i++ {
			wechatUnionMap[list[i].WechatUnionID] = list[i].UserID
		}
		break
	}

	for len(userKeys) > 0 {
		list, err := factory.CustomerRepo.GetCustomersByKey(userKeys)
		if err != nil {
			factory.Worker.Logger().Error(err)
		}
		for i := 0; i < len(list); i++ {
			userKeyMap[list[i].UserKey] = list[i].UserID
		}
		break
	}

	for len(userPhones) > 0 {
		list, err := factory.CustomerRepo.GetCustomersByPhone(userPhones)
		if err != nil {
			factory.Worker.Logger().Error(err)
		}
		for i := 0; i < len(list); i++ {
			userPhoneMap[list[i].Phone] = list[i].UserID
		}
		break
	}

	for _, v := range behaviours {
		customerId, ok := userKeyMap[v.UserKey]
		if ok {
			v.UserId = customerId
			continue
		}

		customerId, ok = wechatUnionMap[v.WechatUnionID]
		if ok {
			v.UserId = customerId
			continue
		}

		customerId, ok = wechatUnionMap[v.UserPhone]
		if ok {
			v.UserId = customerId
			continue
		}
		//暂时不支持临时用户
		// if v.TempUserID == "" {
		// 	continue
		// }
		// v.UserId = factory.CustomerRepo.GetTempUserIDByUUID(v.TempUserID, v.SourceID)
	}
	return
}
