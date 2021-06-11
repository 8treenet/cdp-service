package entity

import (
	"encoding/json"
	"fmt"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

// Customer 客户实体
type Customer struct {
	freedom.Entity
	po.Customer
	Extension        map[string]interface{}
	extensionChanges map[string]interface{}
}

func (entity *Customer) Identity() string {
	return fmt.Sprint(entity.UserID)
}

// SetExtension .
func (entity *Customer) SetExtension(m map[string]interface{}) {
	entity.Extension = m
}

// SetExtend .
func (entity *Customer) GetExtension() map[string]interface{} {
	if entity.Extension != nil {
		return entity.Extension
	}
	return make(map[string]interface{})
}

// GetExtensionChanges .
func (entity *Customer) GetExtensionChanges() map[string]interface{} {
	if entity.extensionChanges == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range entity.extensionChanges {
		result[k] = v
	}
	entity.extensionChanges = nil
	return result
}

func (entity *Customer) UpdateExtensionChanges(putData map[string]interface{}) {
	entity.extensionChanges = map[string]interface{}{}
	for key, v := range putData {
		entity.extensionChanges[key] = v
		entity.Extension[key] = v
	}
}

// MarshalJSON .
func (entity *Customer) MarshalJSON() ([]byte, error) {
	var jsonData struct {
		po.Customer
		Extension map[string]interface{} `json:"extension"`
	}
	jsonData.Customer = entity.Customer
	jsonData.Extension = entity.Extension

	return json.Marshal(jsonData)
}

func (entity *Customer) UpdateByMap(putData map[string]interface{}) error {
	for key, item := range putData {
		i, iErr := utils.ToInt(item)
		switch key {
		case "name":
			entity.SetName(fmt.Sprint(item))
		case "gender":
			if iErr != nil {
				return iErr
			}
			entity.SetGender(i)
		case "age":
			if iErr != nil {
				return iErr
			}
			entity.SetAge(i)
		}
	}
	return nil
}
