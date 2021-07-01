package entity

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

// Customer 客户中介实体
type Intermediary struct {
	freedom.Entity
	Templetes []*po.CustomerExtensionMetadata
}

func (entity *Intermediary) Identity() string {
	return "1001"
}

func (entity *Intermediary) VerifyCustomer(customer *Customer, isNew bool) error {
	mt := map[string]*po.CustomerExtensionMetadata{}

	for _, po := range entity.Templetes {
		mt[po.Variable] = po
		_, ok := customer.Extension[po.Variable]
		if !ok && po.Required == 1 && isNew {
			return fmt.Errorf("缺少必填字段 %s", po.Variable)
		}
	}

	data := customer.Extension
	for key, value := range data {
		po, ok := mt[key]
		if !ok {
			return fmt.Errorf("该字段在模板中不存在 %s", key)
		}

		val := reflect.ValueOf(value)
		switch po.Kind {
		case "String":
			if val.Kind() != reflect.String {
				return fmt.Errorf("错误类型 %v %s:%v", "String", po.Variable, value)
			}
			if po.Reg == "" {
				break
			}

			if ok := regexp.MustCompile(po.Reg).MatchString(value.(string)); !ok {
				return fmt.Errorf("正则匹配失败 %v %s:%v", po.Reg, po.Variable, value)
			}
		case "Float32", "Float64":
			if val.Kind() != reflect.Float32 && val.Kind() != reflect.Float64 {
				return fmt.Errorf("错误类型 %v %s:%v", "Float", po.Variable, value)
			}
		case "DateTime":
			if val.Kind() != reflect.String {
				return fmt.Errorf("错误类型 %v %s:%v", "DateTime", po.Variable, value)
			}
			if _, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(value), time.Local); err != nil {
				return fmt.Errorf("错误类型 %v %s:%v", "DateTime", po.Variable, value)
			}

		case "Date":
			if val.Kind() != reflect.String {
				return fmt.Errorf("错误类型 %v %s:%v", "Date", po.Variable, value)
			}
			if _, err := time.ParseInLocation("2006-01-02", fmt.Sprint(value), time.Local); err != nil {
				return fmt.Errorf("错误类型 %v %s:%v", "Date", po.Variable, value)
			}
		default:
			ok := utils.InSlice([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
				reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
				reflect.Uint64, reflect.Float32, reflect.Float64}, val.Kind())
			if !ok {
				return fmt.Errorf("错误类型 %v %s:%v", "Number", po.Variable, value)
			}
		}
	}

	return nil
}
