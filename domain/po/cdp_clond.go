//Package po generated by 'freedom new-po'
package po

import (
	"time"
)

// Clond .
type Clond struct {
	changes  map[string]interface{}
	ID       int       `gorm:"primaryKey;column:id"`
	Key      string    `gorm:"column:key"`
	Deadline time.Time `gorm:"column:deadline"` // 过期时间
	Created  time.Time `gorm:"column:created"`
	Updated  time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *Clond) TableName() string {
	return "cdp_clond"
}

// Location .
func (obj *Clond) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *Clond) GetChanges() map[string]interface{} {
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
func (obj *Clond) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetKey .
func (obj *Clond) SetKey(key string) {
	obj.Key = key
	obj.Update("key", key)
}

// SetDeadline .
func (obj *Clond) SetDeadline(deadline time.Time) {
	obj.Deadline = deadline
	obj.Update("deadline", deadline)
}

// SetCreated .
func (obj *Clond) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *Clond) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}
