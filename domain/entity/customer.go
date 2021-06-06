package entity

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/8treenet/crm-service/domain/po"
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
func (entity *Customer) Verify() error {
	for _, po := range entity.Templetes {
		value, ok := entity.Source[po.Name]
		if !ok && po.Required == 1 && entity.changes == nil {
			return fmt.Errorf("必填字段 %s", po.Name)
		}
		if !ok {
			continue
		}
		typ := reflect.TypeOf(value)
		switch po.Kind {
		case "String":
			if typ.Kind() != reflect.String {
				return fmt.Errorf("错误类型 %v %s:%v", typ.Kind(), po.Name, value)
			}
			if po.Reg == "" {
				break
			}

			if ok := regexp.MustCompile(po.Reg).MatchString(value.(string)); !ok {
				return fmt.Errorf("正则匹配失败 %v %s:%v", po.Reg, po.Name, value)
			}
		case "Boolean":
			if typ.Kind() != reflect.Bool {
				return fmt.Errorf("错误类型 %v %s:%v", typ.Kind(), po.Name, value)
			}
		case "Double":
			if typ.Kind() != reflect.Float32 && typ.Kind() != reflect.Float64 {
				return fmt.Errorf("错误类型 %v %s:%v", typ.Kind(), po.Name, value)
			}
		default:
			_, err := strconv.Atoi(fmt.Sprint(value))
			if err != nil {
				return fmt.Errorf("错误类型 %v %s:%v", typ.Kind(), po.Name, value)
			}
		}
	}

	return nil
}
