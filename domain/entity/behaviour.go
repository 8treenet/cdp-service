package entity

import (
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Behaviour 客户行为
// 默认ck表字段，customerId, ip，省、市, SouceID
type Behaviour struct {
	freedom.Entity
	po.Behaviour
	CustomerId int    //客户id
	Region     string //省
	City       string //市
}
