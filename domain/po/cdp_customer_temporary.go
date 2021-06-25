//Package po generated by 'freedom new-po'
package po

import (
	"gorm.io/gorm"
	"time"
)

// CustomerTemporary .
type CustomerTemporary struct {
	changes  map[string]interface{}
	ID       int       `gorm:"primaryKey;column:id"`
	UUID     string    `gorm:"column:uuid"` // 临时用户的唯一id
	UserID   string    `gorm:"column:userId"`
	SourceID int       `gorm:"column:sourceId"` // 来源id
	Created  time.Time `gorm:"column:created"`
	Updated  time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *CustomerTemporary) TableName() string {
	return "cdp_customer_temporary"
}

// Location .
func (obj *CustomerTemporary) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *CustomerTemporary) GetChanges() map[string]interface{} {
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
func (obj *CustomerTemporary) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetUUID .
func (obj *CustomerTemporary) SetUUID(uUID string) {
	obj.UUID = uUID
	obj.Update("uuid", uUID)
}

// SetUserID .
func (obj *CustomerTemporary) SetUserID(userID string) {
	obj.UserID = userID
	obj.Update("userId", userID)
}

// SetSourceID .
func (obj *CustomerTemporary) SetSourceID(sourceID int) {
	obj.SourceID = sourceID
	obj.Update("sourceId", sourceID)
}

// SetCreated .
func (obj *CustomerTemporary) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *CustomerTemporary) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddSourceID .
func (obj *CustomerTemporary) AddSourceID(sourceID int) {
	obj.SourceID += sourceID
	obj.Update("sourceId", gorm.Expr("sourceId + ?", sourceID))
}
