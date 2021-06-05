package entity

import (
	"encoding/json"
	"time"

	"github.com/8treenet/crm-service/domain/po"
	"github.com/8treenet/freedom"
	"go.mongodb.org/mongo-driver/bson"
)

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
	return nil
}
