//Package po generated by 'freedom new-po'
package po

import (
	"time"

	"gorm.io/gorm"
)

// CustomerTemplate .
type CustomerTemplate struct {
	changes map[string]interface{}
	ID      int       `gorm:"primaryKey;column:id"`
	Name    string    `gorm:"column:name"`  // key
	Kind    string    `gorm:"column:kind"`  // 类型
	Index   int       `gorm:"column:index"` // 索引0:无,1:btree,2:唯一
	Dict    string    `gorm:"column:dict"`  // 关联字典的key
	Created time.Time `gorm:"column:created"`
	Updated time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *CustomerTemplate) TableName() string {
	return "cdp_customer_template"
}

// Location .
func (obj *CustomerTemplate) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *CustomerTemplate) GetChanges() map[string]interface{} {
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
func (obj *CustomerTemplate) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetName .
func (obj *CustomerTemplate) SetName(name string) {
	obj.Name = name
	obj.Update("name", name)
}

// SetKind .
func (obj *CustomerTemplate) SetKind(kind string) {
	obj.Kind = kind
	obj.Update("kind", kind)
}

// SetIndex .
func (obj *CustomerTemplate) SetIndex(index int) {
	obj.Index = index
	obj.Update("index", index)
}

// SetDict .
func (obj *CustomerTemplate) SetDict(dict string) {
	obj.Dict = dict
	obj.Update("dict", dict)
}

// SetCreated .
func (obj *CustomerTemplate) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *CustomerTemplate) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddIndex .
func (obj *CustomerTemplate) AddIndex(index int) {
	obj.Index += index
	obj.Update("index", gorm.Expr("index + ?", index))
}
