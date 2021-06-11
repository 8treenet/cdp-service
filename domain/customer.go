package domain

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/aggregate"
	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/vo"
	"github.com/8treenet/freedom"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindService(func() *CustomerService {
			return &CustomerService{}
		})
		initiator.InjectController(func(ctx freedom.Context) (service *CustomerService) {
			initiator.FetchService(ctx, &service)
			return
		})
	})
}

// CustomerService .
type CustomerService struct {
	Worker             freedom.Worker
	CustomerRepository *repository.CustomerRepository
	Factory            *aggregate.IntermediaryFactory
}

// GetCustomer 获取客户信息.
func (service *CustomerService) GetCustomer(id int) (*entity.Customer, error) {
	return service.CustomerRepository.GetCustomer(id)
}

// DeleteCustomer 删除客户.
func (service *CustomerService) DeleteCustomer(id int) error {
	customer, err := service.GetCustomer(id)
	if err != nil {
		return err
	}
	return service.CustomerRepository.DeleteCustomer(customer)
}

// UpdateCustomer 修改客户.
func (service *CustomerService) UpdateCustomer(id int, updateOpt map[string]interface{}) error {
	cmd, err := service.Factory.UpdateCustomerNewCmd()
	if err != nil {
		return nil
	}
	return cmd.Do(id, updateOpt)
}

// CreateCustomer 新增客户.
func (service *CustomerService) CreateCustomer(source vo.CustomerDTO) error {
	cmd, err := service.Factory.CreateCustomerNewCmd()
	if err != nil {
		return nil
	}
	return cmd.Do(source)
}

// CreateCustomer 新增客户.
func (service *CustomerService) CreateCustomers(source []vo.CustomerDTO) error {
	cmd, err := service.Factory.CreateCustomerNewCmd()
	if err != nil {
		return nil
	}
	return cmd.BatcheDo(source)
}
