//Package po generated by 'freedom new-po'
package po

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// Behaviour .
type Behaviour struct {
	changes       map[string]interface{}
	ID            int            `gorm:"primaryKey;column:id"`
	WechatUnionID string         `gorm:"column:wechatUnionId"` // 微信唯一id
	UserKey       string         `gorm:"column:userKey"`       // 用户自定义key
	UserPhone     string         `gorm:"column:userPhone"`     // 用户手机号
	TempUserID    string         `gorm:"column:tempUserId"`    // 临时用户唯一id
	UserIPAddr    string         `gorm:"column:userIpAddr"`    // 用户ip地址
	FeatureID     int            `gorm:"column:featureId"`     // 行为的类型
	CreateTime    time.Time      `gorm:"column:createTime"`    // 行为的时间
	Data          datatypes.JSON `gorm:"column:data"`          // 数据
	Processed     int            `gorm:"column:processed"`     // 非0已处理
	SouceID       int            `gorm:"column:souceId"`       // 来源
	Created       time.Time      `gorm:"column:created"`
}

// TableName .
func (obj *Behaviour) TableName() string {
	return "cdp_behaviour"
}

// Location .
func (obj *Behaviour) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *Behaviour) GetChanges() map[string]interface{} {
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
func (obj *Behaviour) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetWechatUnionID .
func (obj *Behaviour) SetWechatUnionID(wechatUnionID string) {
	obj.WechatUnionID = wechatUnionID
	obj.Update("wechatUnionId", wechatUnionID)
}

// SetUserKey .
func (obj *Behaviour) SetUserKey(userKey string) {
	obj.UserKey = userKey
	obj.Update("userKey", userKey)
}

// SetUserPhone .
func (obj *Behaviour) SetUserPhone(userPhone string) {
	obj.UserPhone = userPhone
	obj.Update("userPhone", userPhone)
}

// SetTempUserID .
func (obj *Behaviour) SetTempUserID(tempUserID string) {
	obj.TempUserID = tempUserID
	obj.Update("tempUserId", tempUserID)
}

// SetUserIPAddr .
func (obj *Behaviour) SetUserIPAddr(userIPAddr string) {
	obj.UserIPAddr = userIPAddr
	obj.Update("userIpAddr", userIPAddr)
}

// SetFeatureID .
func (obj *Behaviour) SetFeatureID(featureID int) {
	obj.FeatureID = featureID
	obj.Update("featureId", featureID)
}

// SetCreateTime .
func (obj *Behaviour) SetCreateTime(createTime time.Time) {
	obj.CreateTime = createTime
	obj.Update("createTime", createTime)
}

// SetData .
func (obj *Behaviour) SetData(data datatypes.JSON) {
	obj.Data = data
	obj.Update("data", data)
}

// SetProcessed .
func (obj *Behaviour) SetProcessed(processed int) {
	obj.Processed = processed
	obj.Update("processed", processed)
}

// SetSouceID .
func (obj *Behaviour) SetSouceID(souceID int) {
	obj.SouceID = souceID
	obj.Update("souceId", souceID)
}

// SetCreated .
func (obj *Behaviour) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// AddFeatureID .
func (obj *Behaviour) AddFeatureID(featureID int) {
	obj.FeatureID += featureID
	obj.Update("featureId", gorm.Expr("featureId + ?", featureID))
}

// AddProcessed .
func (obj *Behaviour) AddProcessed(processed int) {
	obj.Processed += processed
	obj.Update("processed", gorm.Expr("processed + ?", processed))
}

// AddSouceID .
func (obj *Behaviour) AddSouceID(souceID int) {
	obj.SouceID += souceID
	obj.Update("souceId", gorm.Expr("souceId + ?", souceID))
}
