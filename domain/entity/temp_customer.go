package entity

import (
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// TempCustomer 临时客户实体
type TempCustomer struct {
	freedom.Entity
	po.CustomerTemporary
}

func (entity *TempCustomer) Identity() string {
	return entity.UserID
}
