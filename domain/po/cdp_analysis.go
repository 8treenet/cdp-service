//Package po generated by 'freedom new-po'
package po

import (
	"gorm.io/gorm"
	"time"
)

// Analysis .
type Analysis struct {
	changes          map[string]interface{}
	ID               int       `gorm:"primaryKey;column:id"`
	Name             string    `gorm:"column:name"`             // 名称
	Title            string    `gorm:"column:title"`            // 显示名称
	FeatureID        int       `gorm:"column:featureId"`        // 元id
	DateRange        int       `gorm:"column:dateRange"`        // 天/范围
	DateConservation int       `gorm:"column:dateConservation"` // 默认0自然日
	DenominatorID    int       `gorm:"column:denominatorId"`    // 分母的统计表id
	OutType          string    `gorm:"column:outType"`          // 输出类型
	XMLData          string    `gorm:"column:xmlData"`          // xml条件数据
	Created          time.Time `gorm:"column:created"`
	Updated          time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *Analysis) TableName() string {
	return "cdp_analysis"
}

// Location .
func (obj *Analysis) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *Analysis) GetChanges() map[string]interface{} {
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
func (obj *Analysis) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetName .
func (obj *Analysis) SetName(name string) {
	obj.Name = name
	obj.Update("name", name)
}

// SetTitle .
func (obj *Analysis) SetTitle(title string) {
	obj.Title = title
	obj.Update("title", title)
}

// SetFeatureID .
func (obj *Analysis) SetFeatureID(featureID int) {
	obj.FeatureID = featureID
	obj.Update("featureId", featureID)
}

// SetDateRange .
func (obj *Analysis) SetDateRange(dateRange int) {
	obj.DateRange = dateRange
	obj.Update("dateRange", dateRange)
}

// SetDateConservation .
func (obj *Analysis) SetDateConservation(dateConservation int) {
	obj.DateConservation = dateConservation
	obj.Update("dateConservation", dateConservation)
}

// SetDenominatorID .
func (obj *Analysis) SetDenominatorID(denominatorID int) {
	obj.DenominatorID = denominatorID
	obj.Update("denominatorId", denominatorID)
}

// SetOutType .
func (obj *Analysis) SetOutType(outType string) {
	obj.OutType = outType
	obj.Update("outType", outType)
}

// SetXMLData .
func (obj *Analysis) SetXMLData(xMLData string) {
	obj.XMLData = xMLData
	obj.Update("xmlData", xMLData)
}

// SetCreated .
func (obj *Analysis) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *Analysis) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddFeatureID .
func (obj *Analysis) AddFeatureID(featureID int) {
	obj.FeatureID += featureID
	obj.Update("featureId", gorm.Expr("featureId + ?", featureID))
}

// AddDateRange .
func (obj *Analysis) AddDateRange(dateRange int) {
	obj.DateRange += dateRange
	obj.Update("dateRange", gorm.Expr("dateRange + ?", dateRange))
}

// AddDateConservation .
func (obj *Analysis) AddDateConservation(dateConservation int) {
	obj.DateConservation += dateConservation
	obj.Update("dateConservation", gorm.Expr("dateConservation + ?", dateConservation))
}

// AddDenominatorID .
func (obj *Analysis) AddDenominatorID(denominatorID int) {
	obj.DenominatorID += denominatorID
	obj.Update("denominatorId", gorm.Expr("denominatorId + ?", denominatorID))
}
