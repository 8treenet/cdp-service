package domain

import (
	"encoding/json"
	"net"
	"time"

	"cdp-service/adapter/repository"
	"cdp-service/domain/aggregate"
	"cdp-service/domain/po"
	"cdp-service/domain/vo"
	"cdp-service/utils"

	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *BehaviourService {
			return &BehaviourService{fetchTime: 3 * time.Second, fetchCount: 800}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *BehaviourService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// BehaviourService .
type BehaviourService struct {
	Worker              freedom.Worker
	BehaviourRepository *repository.BehaviourRepository
	SupportRepository   *repository.SupportRepository
	FeatureRepository   *repository.FeatureRepository
	CustomerRepository  *repository.CustomerRepository
	BehaviourFactory    *aggregate.BehaviourFactory
	ClondRepository     *repository.ClondRepository
	fetchTime           time.Duration
	fetchCount          int
}

// CreateBehaviour
func (service *BehaviourService) CreateBehaviours(featureID int, reqs []vo.ReqBehaviourDTO) error {
	fentity, err := service.FeatureRepository.GetFeatureEntity(featureID)
	if err != nil {
		return err
	}

	for _, req := range reqs {
		if err := fentity.CheckMetadata(req.Data); err != nil {
			return err
		}
		jsonData, err := json.Marshal(req.Data)
		if err != nil {
			return err
		}

		sourceId := service.SupportRepository.FindSourceID(req.Source)
		createTime := time.Now()
		inCreateTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.CreateTime, time.Local)
		if err == nil {
			createTime = inCreateTime
		}
		ipAddr := ""
		if net.ParseIP(req.IPAddr) != nil {
			ipAddr = req.IPAddr
		}

		obj := &po.Behaviour{
			WechatUnionID: req.WechatUnionID,
			UserKey:       req.UserKey,
			UserPhone:     req.UserPhone,
			UserIPAddr:    ipAddr,
			FeatureID:     fentity.ID,
			CreateTime:    createTime,
			Data:          jsonData,
			SourceID:      sourceId,
			Created:       time.Now(),
		}
		service.BehaviourRepository.AddQueue([]*po.Behaviour{obj})
	}
	return nil
}

// CreateBehavioursByCSV
func (service *BehaviourService) CreateBehavioursByCSV(featureID int, key string) error {
	data, err := service.ClondRepository.PublicDownload(key)
	if err != nil {
		return err
	}
	fentity, err := service.FeatureRepository.GetFeatureEntity(featureID)
	if err != nil {
		return err
	}

	csv, err := utils.NewCSV(data)
	if err != nil {
		return err
	}

	var pos []*po.Behaviour
	for _, v := range csv.ToMaps() {
		metadata, err := fentity.ConvertMetadata(v)
		if err != nil {
			return err
		}
		jsonData, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		createTime := time.Now()
		inCreateTime, err := time.ParseInLocation("2006-01-02 15:04:05", v["createTime"], time.Local)
		if err == nil {
			createTime = inCreateTime
		}

		sourceId := service.SupportRepository.FindSourceID(v["source"])
		ipAddr := ""
		if net.ParseIP(v["ipAddr"]) != nil {
			ipAddr = v["ipAddr"]
		}

		obj := &po.Behaviour{
			WechatUnionID: v["wechatUnionID"],
			UserKey:       v["userKey"],
			UserPhone:     v["userPhone"],
			FeatureID:     fentity.ID,
			CreateTime:    createTime,
			Data:          jsonData,

			SourceID:   sourceId,
			UserIPAddr: ipAddr,
			Created:    time.Now(),
		}
		pos = append(pos, obj)
	}

	service.BehaviourRepository.AddQueue(pos)
	return nil
}

// BatchSave 批量入库
func (service *BehaviourService) BatchSave() func() bool {
	list, cancel := service.BehaviourRepository.FetchQueue(service.fetchCount, service.fetchTime)

	if len(list) == 0 {
		return cancel
	}
	//i<2重试一次
	for i := 0; i < 2; i++ {
		err := service.batchTempCustomer(list)
		if err == nil {
			break
		}
		if err != nil && i == 1 {
			service.Worker.Logger().Error(err)
		}
		time.Sleep(3 * time.Second)
	}

	//i<2重试一次
	for i := 0; i < 2; i++ {
		err := service.BehaviourRepository.BatchSave(list)
		if err == nil {
			break
		}
		if err != nil && i == 1 {
			service.Worker.Logger().Error(err)
		}
		time.Sleep(3 * time.Second)
	}
	return cancel
}

// EnteringHouse 入数仓
func (service *BehaviourService) EnteringHouse() {
	cmds, err := service.BehaviourFactory.CreateBehaviourCmds()
	if err != nil {
		service.Worker.Logger().Error(err)
		return
	}

	for _, cmd := range cmds {
		if err := cmd.Do(); err != nil {
			service.Worker.Logger().Error(err)
		}
	}
}

// Truncate .
func (service *BehaviourService) Truncate() {
	for i := 0; i < 120; i++ {
		ok := service.BehaviourRepository.TruncateBehaviour()
		if !ok {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return
}

// batchTempCustomer 处理临时用户
func (service *BehaviourService) batchTempCustomer(list []*po.Behaviour) error {
	tempCustomerMap := map[int][]string{}

	for _, behaviour := range list {
		if behaviour.TempUserID == "" {
			continue
		}

		_, ok := tempCustomerMap[behaviour.SourceID]
		if !ok {
			tempCustomerMap[behaviour.SourceID] = append(tempCustomerMap[behaviour.SourceID], behaviour.TempUserID)
			continue
		}
		if utils.InSlice(tempCustomerMap[behaviour.SourceID], behaviour.TempUserID) {
			continue
		}
		tempCustomerMap[behaviour.SourceID] = append(tempCustomerMap[behaviour.SourceID], behaviour.TempUserID)
	}

	for sourceId, ids := range tempCustomerMap {
		existIds := service.CustomerRepository.GetExistTempCustomers(ids, sourceId) //已存在的临时用户
		newIds := []string{}
		for _, id := range ids {
			if utils.InSlice(existIds, id) {
				continue
			}
			newIds = append(newIds, id)
		}
		if len(newIds) == 0 {
			continue
		}

		newTempCustomers := []*po.CustomerTemporary{}
		for _, id := range newIds {
			userId, _ := utils.GenerateUUID()
			newTempCustomers = append(newTempCustomers, &po.CustomerTemporary{
				UUID:     id,
				UserID:   userId,
				SourceID: sourceId,
				Created:  time.Now(),
				Updated:  time.Now(),
			})
		}

		if err := service.CustomerRepository.CreateTempCustomer(newTempCustomers); err != nil {
			service.Worker.Logger().Error(err)
			return err
		}
	}

	return nil
}
