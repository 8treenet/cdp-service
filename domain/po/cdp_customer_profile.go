//Package po generated by 'freedom new-po'
package po

import (
	"time"
)

// CustomerProfile .
type CustomerProfile struct {
	changes   map[string]interface{}
	UserID    string    `gorm:"primaryKey;column:userId"`
	PersonaID int       `gorm:"primaryKey;column:personaId"` // 画像id
	BeginTime time.Time `gorm:"column:beginTime"`
	EndTime   time.Time `gorm:"column:endTime"`
	Created   time.Time `gorm:"column:created"`
	Updated   time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *CustomerProfile) TableName() string {
	return "cdp_customer_profile"
}

// Location .
func (obj *CustomerProfile) Location() map[string]interface{} {
	return map[string]interface{}{"personaId": obj.PersonaID, "userId": obj.UserID}
}

// GetChanges .
func (obj *CustomerProfile) GetChanges() map[string]interface{} {
	if obj.changes == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range obj.changes {
		result[k] = v
	}
	obj.changes = nil
	return result
}

// Update .
func (obj *CustomerProfile) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetBeginTime .
func (obj *CustomerProfile) SetBeginTime(beginTime time.Time) {
	obj.BeginTime = beginTime
	obj.Update("beginTime", beginTime)
}

// SetEndTime .
func (obj *CustomerProfile) SetEndTime(endTime time.Time) {
	obj.EndTime = endTime
	obj.Update("endTime", endTime)
}

// SetCreated .
func (obj *CustomerProfile) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *CustomerProfile) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}
