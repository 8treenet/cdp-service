package entity

import (
	"encoding/json"
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

func (entity *Feature) AddMetadata(variable, title, kind, dict string) {
	entity.FeatureMetadata = append(entity.FeatureMetadata, &po.BehaviourFeatureMetadata{
		Variable: variable,
		Title:    title,
		Kind:     kind,
		Dict:     dict,
		Created:  time.Now(),
		Updated:  time.Now(),
	})
}

// MarshalJSON .
func (entity *Feature) MarshalJSON() ([]byte, error) {
	var jsonData struct {
		ID        int         `json:"id"`
		Title     string      `json:"title"`
		Warehouse string      `json:"warehouse"`
		Created   string      `json:"created"`
		Metadata  interface{} `json:"metadata"`
	}

	jsonData.ID = entity.ID
	jsonData.Title = entity.Title
	jsonData.Warehouse = entity.Warehouse
	jsonData.Created = entity.Created.Format("2006-01-02 15:04:05")

	list := make([]struct {
		Variable string `json:"variable"`
		Title    string `json:"title"`
		Kind     string `json:"kind"`
		Dict     string `json:"dict"`
	}, len(entity.FeatureMetadata))

	for i := 0; i < len(entity.FeatureMetadata); i++ {
		list[i].Variable = entity.FeatureMetadata[i].Variable
		list[i].Title = entity.FeatureMetadata[i].Title
		list[i].Kind = entity.FeatureMetadata[i].Kind
		list[i].Dict = entity.FeatureMetadata[i].Dict
	}
	jsonData.Metadata = list

	return json.Marshal(jsonData)
}
