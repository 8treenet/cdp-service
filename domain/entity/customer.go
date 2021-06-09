package entity

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
	"go.mongodb.org/mongo-driver/bson"
)

// Customer 客户实体
type Customer struct {
	freedom.Entity
	Source    map[string]interface{}
	Templetes []*po.CustomerTemplate

	changes map[string]interface{}
}

func (entity *Customer) Identity() string {
	if entity.Source == nil {
		return ""
	}

	iid, ok := entity.Source["_id"]
	if !ok {
		return ""
	}
	id, _ := iid.(string)
	return id
}

// Location .
func (entity *Customer) Location() interface{} {
	return bson.D{{"_id", entity.Identity()}}
}

// GetChanges .
func (entity *Customer) GetChanges() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range entity.changes {
		result[k] = v
	}
	result["_updated"] = time.Now()
	entity.changes = nil
	return result
}

// Update .
func (entity *Customer) Update(name string, value interface{}) {
	if entity.changes == nil {
		entity.changes = make(map[string]interface{})
	}
	entity.changes[name] = value
	entity.Source[name] = value
}

// MarshalJSON .
func (entity *Customer) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{"_id": entity.Identity()}

	for _, po := range entity.Templetes {
		value, ok := entity.Source[po.Name]
		if !ok {
			continue
		}
		data[po.Name] = value
	}
	return json.Marshal(data)
}

// Verify .
func (entity *Customer) Verify(isNew ...bool) error {
	mt := map[string]*po.CustomerTemplate{}
	for _, po := range entity.Templetes {
		mt[po.Name] = po
		_, ok := entity.Source[po.Name]
		if !ok && po.Required == 1 && len(isNew) > 0 {
			return fmt.Errorf("缺少必填字段 %s", po.Name)
		}
	}
	data := entity.Source
	if len(isNew) == 0 {
		data = entity.changes
	}

	for key, value := range data {
		po, ok := mt[key]
		if utils.InSlice([]string{"_id", "_updated", "_created"}, key) {
			continue
		}
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
