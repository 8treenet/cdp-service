package entity

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Customer 客户实体
type Customer struct {
	freedom.Entity
	po.Customer
	Extension map[string]interface{}
}

func (entity *Customer) Identity() string {
	return entity.UserID
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

func (entity *Customer) updateExtension(putData map[string]interface{}) {
	for key, v := range putData {
		entity.Extension[key] = v
	}
}

func (entity *Customer) UpdateByMap(putData map[string]interface{}) error {
	i, ok := putData["extension"]
	if ok {
		extensionMap, _ := i.(map[string]interface{})
		entity.updateExtension(extensionMap)
	}
	for key, item := range putData {
		switch key {
		case "name":
			entity.SetName(fmt.Sprint(item))
		case "gender":
			entity.SetGender(fmt.Sprint(item))
		case "email":
			entity.SetEmail(fmt.Sprint(item))
		case "birthday":
			_, err := time.Parse("2006-01-02", fmt.Sprint(item))
			if err != nil {
				return err
			}
			entity.SetBirthday(fmt.Sprint(item))
		case "province":
			entity.SetProvince(fmt.Sprint(item))
		case "city":
			entity.SetCity(fmt.Sprint(item))
		case "region":
			entity.SetRegion(fmt.Sprint(item))
		}
	}
	return nil
}
