package domain

import (
	"cdp-service/adapter/repository"
	"cdp-service/domain/aggregate"
	"cdp-service/domain/entity"
	"cdp-service/domain/vo"

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
func (service *CustomerService) GetCustomer(id string) (*entity.Customer, error) {
	return service.CustomerRepository.GetCustomer(id)
}

// DeleteCustomer 删除客户.
func (service *CustomerService) DeleteCustomer(id []string) error {
	for i := 0; i < len(id); i++ {
		customer, err := service.GetCustomer(id[i])
		if err != nil {
			service.Worker.Logger().Error(err)
			continue
		}

		if err = service.CustomerRepository.DeleteCustomer(customer); err != nil {
			return err
		}
	}
	return nil
}

// UpdateCustomer 修改客户.
func (service *CustomerService) UpdateCustomer(id string, updateOpt map[string]interface{}) error {
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

// GetCustomersByPage 获取客户信息列表.
func (service *CustomerService) GetCustomersByPage() (result []*entity.Customer, totalPage int, e error) {
	return service.CustomerRepository.GetCustomersByPage()
}

// GetCustomersByKeys 获取客户信息列表.
func (service *CustomerService) GetCustomersByKeys(keys []string) (result []*entity.Customer, e error) {
	return service.CustomerRepository.GetCustomersByKey(keys)
}

// GetCustomersByPhone 获取客户信息列表.
func (service *CustomerService) GetCustomersByPhone(phones []string) (result []*entity.Customer, e error) {
	return service.CustomerRepository.GetCustomersByPhone(phones)
}

// GetCustomersByWechat 获取客户信息列表.
func (service *CustomerService) GetCustomersByWechat(unionIds []string) (result []*entity.Customer, e error) {
	return service.CustomerRepository.GetCustomersByWechat(unionIds)
}
