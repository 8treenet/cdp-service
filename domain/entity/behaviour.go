package entity

import (
	"encoding/json"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/freedom"
)

// Behaviour 客户行为
// 默认ck表字段，customerId, ip，省、市, SourceID
type Behaviour struct {
	freedom.Entity
	po.Behaviour
	UserId string //客户id
	Region string //省
	City   string //市
}

// ToColumns 返回列和值
func (entity *Behaviour) ToColumns() (result map[string]interface{}, e error) {
	result = make(map[string]interface{})

	result[cattle.ColumnCity] = entity.City
	result[cattle.ColumnRegion] = entity.Region
	result[cattle.ColumnIP] = entity.UserIPAddr
	result[cattle.ColumnSourceId] = entity.SourceID
	result[cattle.ColumnUserId] = entity.UserId
	result[cattle.ColumnCreateTime] = entity.CreateTime

	dataMap := map[string]interface{}{}
	e = json.Unmarshal(entity.Data, &dataMap)
	if e != nil {
		return
	}
	for k, v := range dataMap {
		if _, ok := result[k]; ok {
			continue
		}

		result[k] = v
	}
	return
}

// AsyncFinish 同步完成
func (entity *Behaviour) SyncSuccess() {
	entity.SetProcessed(2)
}

// AsyncFinish 同步失败
func (entity *Behaviour) SyncError() {
	entity.SetProcessed(3)
}
