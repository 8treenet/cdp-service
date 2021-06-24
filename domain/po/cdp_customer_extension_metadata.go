//Package po generated by 'freedom new-po'
package po

import (
	"time"

	"gorm.io/gorm"
)

// CustomerExtensionMetadata .
type CustomerExtensionMetadata struct {
	changes  map[string]interface{}
	ID       int       `gorm:"primaryKey;column:id" json:"id"`
	Variable string    `gorm:"column:variable" json:"variable"` // 类型的名称
	Title    string    `gorm:"column:title" json:"title"`       // 中文的标题
	Kind     string    `gorm:"column:kind" json:"kind"`         // 类型
	Dict     string    `gorm:"column:dict" json:"dict"`         // 关联字典的key
	Reg      string    `gorm:"column:reg" json:"reg"`           // 正则
	Required int       `gorm:"column:required" json:"required"` // 1 必填
	Sort     int       `gorm:"column:sort" json:"sort"`         // 排序
	Created  time.Time `gorm:"column:created" json:"-"`
	Updated  time.Time `gorm:"column:updated" json:"-"`
}

// TableName .
func (obj *CustomerExtensionMetadata) TableName() string {
	return "cdp_customer_extension_metadata"
}

// Location .
func (obj *CustomerExtensionMetadata) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *CustomerExtensionMetadata) GetChanges() map[string]interface{} {
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
func (obj *CustomerExtensionMetadata) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetVariable .
func (obj *CustomerExtensionMetadata) SetVariable(variable string) {
	obj.Variable = variable
	obj.Update("variable", variable)
}

// SetTitle .
func (obj *CustomerExtensionMetadata) SetTitle(title string) {
	obj.Title = title
	obj.Update("title", title)
}

// SetKind .
func (obj *CustomerExtensionMetadata) SetKind(kind string) {
	obj.Kind = kind
	obj.Update("kind", kind)
}

// SetDict .
func (obj *CustomerExtensionMetadata) SetDict(dict string) {
	obj.Dict = dict
	obj.Update("dict", dict)
}

// SetReg .
func (obj *CustomerExtensionMetadata) SetReg(reg string) {
	obj.Reg = reg
	obj.Update("reg", reg)
}

// SetRequired .
func (obj *CustomerExtensionMetadata) SetRequired(required int) {
	obj.Required = required
	obj.Update("required", required)
}

// SetSort .
func (obj *CustomerExtensionMetadata) SetSort(sort int) {
	obj.Sort = sort
	obj.Update("sort", sort)
}

// SetCreated .
func (obj *CustomerExtensionMetadata) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *CustomerExtensionMetadata) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddRequired .
func (obj *CustomerExtensionMetadata) AddRequired(required int) {
	obj.Required += required
	obj.Update("required", gorm.Expr("required + ?", required))
}

// AddSort .
func (obj *CustomerExtensionMetadata) AddSort(sort int) {
	obj.Sort += sort
	obj.Update("sort", gorm.Expr("sort + ?", sort))
}
