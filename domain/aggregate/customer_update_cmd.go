package aggregate

import (
	"cdp-service/adapter/repository"
	"cdp-service/domain/entity"

	"github.com/8treenet/freedom/infra/transaction"
)

// CustomerUpdateCmd
type CustomerUpdateCmd struct {
	entity.Intermediary
	CustomerRepo *repository.CustomerRepository
	TX           transaction.Transaction //依赖倒置事务组件
}

// Do .
func (cmd *CustomerUpdateCmd) Do(id string, m map[string]interface{}) (e error) {
	entity, e := cmd.CustomerRepo.GetCustomer(id)
	if e != nil {
		return e
	}

	if e = entity.UpdateByMap(m); e != nil {
		return
	}
	if e = cmd.VerifyCustomer(entity, false); e != nil {
		return
	}

	return cmd.TX.Execute(func() error {
		return cmd.CustomerRepo.SaveCustomer(entity)
	})
}
