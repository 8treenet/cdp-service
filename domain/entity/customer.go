package entity

import (
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

func (customer *Customer) Identity() string {
	if customer.Source == nil {
		return ""
	}

	iid, ok := customer.Source["_id"]
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
func (obj *Customer) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
	obj.Source[name] = value
}
