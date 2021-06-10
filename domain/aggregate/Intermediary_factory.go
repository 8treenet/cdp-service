package aggregate

import (
	"github.com/8treenet/cdp-service/adapter/repository"
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
	CustomerRepo *repository.CustomerRepository
	Intermediary *repository.IntermediaryRepository
	TX           transaction.Transaction //依赖倒置事务组件
	Worker       freedom.Worker          //运行时，一个请求绑定一个运行时
}

// CreateCustomerNewCmd 返回添加客户命令
func (factory *IntermediaryFactory) CreateCustomerNewCmd() (cmd *CustomerNewCmd, e error) {
	ientity, err := factory.Intermediary.GetEntity()
	if err != nil {
		return nil, err
	}

	return &CustomerNewCmd{
		Intermediary: *ientity,
		CustomerRepo: factory.CustomerRepo,
		TX:           factory.TX,
	}, nil
}
