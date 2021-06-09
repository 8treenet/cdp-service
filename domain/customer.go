package domain

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/domain/entity"
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
}

// GetCustomer 获取客户信息.
func (service *CustomerService) GetCustomer(id string) (*entity.Customer, error) {
	return service.CustomerRepository.GetCustomer(id)
}

// DeleteCustomer 删除客户.
func (service *CustomerService) DeleteCustomer(id string) error {
	customer, err := service.GetCustomer(id)
	if err != nil {
		return err
	}
	return service.CustomerRepository.DeleteCustomer(customer)
}

// UpdateCustomer 修改客户.
func (service *CustomerService) UpdateCustomer(id string, updateOpt map[string]interface{}) error {
	customer, err := service.GetCustomer(id)
	if err != nil {
		return err
	}

	for name, v := range updateOpt {
		customer.Update(name, v)
	}
	return service.CustomerRepository.SaveCustomer(customer)
}

// CreateCustomer 新增客户.
func (service *CustomerService) CreateCustomer(source map[string]interface{}) error {
	_, err := service.CustomerRepository.NewCustomer(source)
	return err
}

// CreateCustomers 新增客户.
func (service *CustomerService) CreateCustomers(sources []map[string]interface{}) error {
	_, err := service.CustomerRepository.NewCustomers(sources)
	return err
}
