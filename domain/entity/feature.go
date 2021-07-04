package entity

import (
	"time"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/freedom"
)

// Feature 特征实体
type Feature struct {
	po.BehaviourFeature
	freedom.Entity
	FeatureMetadata []*po.BehaviourFeatureMetadata
}

func (entity *Feature) Identity() string {
	return "1001"
}

func (entity *Feature) AddMetadata(variable, title, kind, dict string, orderByNumber, partition int) {
	entity.FeatureMetadata = append(entity.FeatureMetadata, &po.BehaviourFeatureMetadata{
		Variable:      variable,
		Title:         title,
		Kind:          kind,
		Dict:          dict,
		OrderByNumber: orderByNumber,
		Partition:     partition,
		Created:       time.Now(),
		Updated:       time.Now(),
	})
}

func (entity *Feature) View() interface{} {
	var jsonData struct {
		ID           int         `json:"id"`
		Title        string      `json:"title"`
		Warehouse    string      `json:"warehouse"`
		Created      string      `json:"created"`
		CategoryType int         `json:"categoryType"` // 0自定义行为，1系统提供行为，2系统提供不可扩展
		Category     string      `json:"category"`     // 行业
		Metadata     interface{} `json:"metadata"`
	}

	jsonData.ID = entity.ID
	jsonData.Title = entity.Title
	jsonData.Warehouse = entity.Warehouse
	jsonData.Created = entity.Created.Format("2006-01-02 15:04:05")
	jsonData.CategoryType = entity.CategoryType
	jsonData.Category = entity.Category

	list := make([]struct {
		Variable      string `json:"variable"`
		Title         string `json:"title"`
		Kind          string `json:"kind"`
		Dict          string `json:"dict"`
		OrderByNumber int    `json:"orderByNumber"` // ck排序键，非0排序
		Partition     int    `json:"partition"`     // 非0分区
	}, len(entity.FeatureMetadata))

	for i := 0; i < len(entity.FeatureMetadata); i++ {
		list[i].Variable = entity.FeatureMetadata[i].Variable
		list[i].Title = entity.FeatureMetadata[i].Title
		list[i].Kind = entity.FeatureMetadata[i].Kind
		list[i].Dict = entity.FeatureMetadata[i].Dict
		list[i].OrderByNumber = entity.FeatureMetadata[i].OrderByNumber
		list[i].Partition = entity.FeatureMetadata[i].Partition
	}
	jsonData.Metadata = list
	return jsonData
}

// ToColumns 返回列和类型
func (entity *Feature) ToColumns() (result map[string]string) {
	result = make(map[string]string)
	for _, v := range entity.FeatureMetadata {
		result[v.Variable] = v.Kind
	}
	return result
}
