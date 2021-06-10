package aggregate

import (
	"time"

	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom/infra/transaction"
)

// CustomerNewCmd
type CustomerNewCmd struct {
	entity.Intermediary
	CustomerRepo *repository.CustomerRepository
	TX           transaction.Transaction //依赖倒置事务组件
}

// Do .
func (cmd *CustomerNewCmd) Do(customerDto vo.CustomerDTO) (e error) {
	customer := cmd.CustomerRepo.NewCustomer()
	customer.Customer = customerDto.Customer
	customer.Customer.UserID = 0
	customer.Customer.Created = time.Now()
	customer.Customer.Updated = time.Now()

	customer.SetExtend(customerDto.Extend)
	if e = cmd.VerifyCustomer(customer, true); e != nil {
		return
	}

	return cmd.TX.Execute(func() error {
		return cmd.CustomerRepo.SaveCustomer(customer)
	})
}

// BatcheDo .
func (cmd *CustomerNewCmd) BatcheDo(customerDtos []vo.CustomerDTO) (e error) {
	for _, v := range customerDtos {
		if e = cmd.Do(v); e != nil {
			return
		}
	}
	return
}
