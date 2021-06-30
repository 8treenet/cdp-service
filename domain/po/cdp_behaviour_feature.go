//Package po generated by 'freedom new-po'
package po

import (
	"time"

	"gorm.io/gorm"
)

// BehaviourFeature .
type BehaviourFeature struct {
	changes      map[string]interface{}
	ID           int       `gorm:"primaryKey;column:id"`
	Title        string    `gorm:"column:title"`
	Warehouse    string    `gorm:"column:warehouse"`    // clickhouse的表名
	CategoryType int       `gorm:"column:categoryType"` // 0自定义行为，1系统提供行为，2系统提供不可扩展
	Category     string    `gorm:"column:category"`     // 行业
	Created      time.Time `gorm:"column:created"`
	Updated      time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *BehaviourFeature) TableName() string {
	return "cdp_behaviour_feature"
}

// Location .
func (obj *BehaviourFeature) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *BehaviourFeature) GetChanges() map[string]interface{} {
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
func (obj *BehaviourFeature) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetTitle .
func (obj *BehaviourFeature) SetTitle(title string) {
	obj.Title = title
	obj.Update("title", title)
}

// SetWarehouse .
func (obj *BehaviourFeature) SetWarehouse(warehouse string) {
	obj.Warehouse = warehouse
	obj.Update("warehouse", warehouse)
}

// SetCategoryType .
func (obj *BehaviourFeature) SetCategoryType(categoryType int) {
	obj.CategoryType = categoryType
	obj.Update("categoryType", categoryType)
}

// SetCategory .
func (obj *BehaviourFeature) SetCategory(category string) {
	obj.Category = category
	obj.Update("category", category)
}

// SetCreated .
func (obj *BehaviourFeature) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *BehaviourFeature) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddCategoryType .
func (obj *BehaviourFeature) AddCategoryType(categoryType int) {
	obj.CategoryType += categoryType
	obj.Update("categoryType", gorm.Expr("categoryType + ?", categoryType))
}
