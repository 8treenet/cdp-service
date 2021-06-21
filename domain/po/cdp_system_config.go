//Package po generated by 'freedom new-po'
package po

import (
	"time"
)

// SystemConfig .
type SystemConfig struct {
	changes map[string]interface{}
	ID      int       `gorm:"primaryKey;column:id"`
	Name    string    `gorm:"column:name"`  // 配置名称
	Value   string    `gorm:"column:value"` // 配置数据
	Created time.Time `gorm:"column:created"`
	Updated time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *SystemConfig) TableName() string {
	return "cdp_system_config"
}

// Location .
func (obj *SystemConfig) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *SystemConfig) GetChanges() map[string]interface{} {
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
func (obj *SystemConfig) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetName .
func (obj *SystemConfig) SetName(name string) {
	obj.Name = name
	obj.Update("name", name)
}

// SetValue .
func (obj *SystemConfig) SetValue(value string) {
	obj.Value = value
	obj.Update("value", value)
}

// SetCreated .
func (obj *SystemConfig) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *SystemConfig) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}
