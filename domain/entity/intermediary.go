package entity

import (
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Customer 客户中介实体
type Intermediary struct {
	freedom.Entity
	Templetes []*po.CustomerExtendTemplate
}

func (entity *Intermediary) VerifyCustomer(customer *Customer, new bool) error {
	return nil
}

func (entity *Intermediary) Identity() string {
	return "1001"
}
