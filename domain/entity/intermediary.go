package entity

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

// Customer 客户中介实体
type Intermediary struct {
	freedom.Entity
	Templetes []*po.CustomerExtensionTemplate
}

func (entity *Intermediary) Identity() string {
	return "1001"
}

func (entity *Intermediary) VerifyCustomer(customer *Customer, isNew bool) error {
	mt := map[string]*po.CustomerExtensionTemplate{}

	for _, po := range entity.Templetes {
		mt[po.Name] = po
		_, ok := customer.Extension[po.Name]
		if !ok && po.Required == 1 && isNew {
			return fmt.Errorf("缺少必填字段 %s", po.Name)
		}
	}

	data := customer.Extension
	if !isNew {
		data = customer.extensionChanges
	}

	for key, value := range data {
		po, ok := mt[key]
		if !ok {
			return fmt.Errorf("该字段在模板中不存在 %s", key)
		}

		val := reflect.ValueOf(value)
		switch po.Kind {
		case "String":
			if val.Kind() != reflect.String {
				return fmt.Errorf("错误类型 %v %s:%v", "String", po.Name, value)
			}
			if po.Reg == "" {
				break
			}

			if ok := regexp.MustCompile(po.Reg).MatchString(value.(string)); !ok {
				return fmt.Errorf("正则匹配失败 %v %s:%v", po.Reg, po.Name, value)
			}
		case "Boolean":
			if val.Kind() != reflect.Bool {
				return fmt.Errorf("错误类型 %v %s:%v", "Boolean", po.Name, value)
			}
		case "Double":
			if val.Kind() != reflect.Float32 && val.Kind() != reflect.Float64 {
				return fmt.Errorf("错误类型 %v %s:%v", "Double", po.Name, value)
			}
		default:
			ok := utils.InSlice([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
				reflect.Uint64, reflect.Float32, reflect.Float64}, val.Kind())
			if !ok {
				return fmt.Errorf("错误类型 %v %s:%v", "Integer", po.Name, value)
			}
		}
	}

	return nil
}
