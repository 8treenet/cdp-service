package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/infra/transaction"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindFactory(func() *IntermediaryFactory {
			return &IntermediaryFactory{} //创建中介工厂
		})
	})
}

// IntermediaryFactory 中介工厂
type IntermediaryFactory struct {
	CustomerRepo        *repository.CustomerRepository
	Intermediary        *repository.IntermediaryRepository
	SignRepo            *repository.SignRepository
	SupportRepository   *repository.SupportRepository
	TX                  transaction.Transaction //依赖倒置事务组件
	Worker              freedom.Worker          //运行时，一个请求绑定一个运行时
	GEO                 *infra.GEO              //geo
	BehaviourRepository *repository.BehaviourRepository
}

// CreateCustomerNewCmd 返回添加客户命令
func (factory *IntermediaryFactory) CreateCustomerNewCmd() (cmd *CustomerCreateCmd, e error) {
	ientity, err := factory.Intermediary.GetEntity()
	if err != nil {
		return nil, err
	}

	userRegisterEntity, err := factory.SupportRepository.GetFeatureEntityByWarehouse("user_register")
	if err != nil {
		return nil, err
	}

	return &CustomerCreateCmd{
		Intermediary:        *ientity,
		CustomerRepo:        factory.CustomerRepo,
		SignRepo:            factory.SignRepo,
		TX:                  factory.TX,
		SupportRepository:   factory.SupportRepository,
		GEO:                 factory.GEO,
		BehaviourRepository: factory.BehaviourRepository,
		UserRegisterEntity:  userRegisterEntity,
	}, nil
}

// CreateCustomerNewCmd 返回添加客户命令
func (factory *IntermediaryFactory) UpdateCustomerNewCmd() (cmd *CustomerUpdateCmd, e error) {
	ientity, err := factory.Intermediary.GetEntity()
	if err != nil {
		return nil, err
	}

	return &CustomerUpdateCmd{
		Intermediary: *ientity,
		CustomerRepo: factory.CustomerRepo,
		TX:           factory.TX,
	}, nil
}
