//Package po generated by 'freedom new-po'
package po

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// AnalysisReport .
type AnalysisReport struct {
	changes    map[string]interface{}
	ID         int            `gorm:"primaryKey;column:id"`
	AnalysisID int            `gorm:"column:analysisId"` // id
	Data       datatypes.JSON `gorm:"column:data"`       // 结果
	BeginTime  time.Time      `gorm:"column:beginTime"`
	EndTime    time.Time      `gorm:"column:endTime"`
	Created    time.Time      `gorm:"column:created"`
	Updated    time.Time      `gorm:"column:updated"`
}

// TableName .
func (obj *AnalysisReport) TableName() string {
	return "cdp_analysis_report"
}

// Location .
func (obj *AnalysisReport) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *AnalysisReport) GetChanges() map[string]interface{} {
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
func (obj *AnalysisReport) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetAnalysisID .
func (obj *AnalysisReport) SetAnalysisID(analysisID int) {
	obj.AnalysisID = analysisID
	obj.Update("analysisId", analysisID)
}

// SetData .
func (obj *AnalysisReport) SetData(data datatypes.JSON) {
	obj.Data = data
	obj.Update("data", data)
}

// SetBeginTime .
func (obj *AnalysisReport) SetBeginTime(beginTime time.Time) {
	obj.BeginTime = beginTime
	obj.Update("beginTime", beginTime)
}

// SetEndTime .
func (obj *AnalysisReport) SetEndTime(endTime time.Time) {
	obj.EndTime = endTime
	obj.Update("endTime", endTime)
}

// SetCreated .
func (obj *AnalysisReport) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *AnalysisReport) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddAnalysisID .
func (obj *AnalysisReport) AddAnalysisID(analysisID int) {
	obj.AnalysisID += analysisID
	obj.Update("analysisId", gorm.Expr("analysisId + ?", analysisID))
}
