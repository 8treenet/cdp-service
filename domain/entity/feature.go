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

// MarshalJSON .
func (entity *Feature) MarshalJSON() ([]byte, error) {
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

	return json.Marshal(jsonData)
}

// VerifyBehaviour 验证数据
func (entity *Feature) VerifyBehaviour(customer *Behaviour) error {
	//mt := map[string]*po.BehaviourFeatureMetadata{}

	// for _, po := range entity.Templetes {
	// 	mt[po.Variable] = po
	// 	_, ok := customer.Extension[po.Variable]
	// 	if !ok && po.Required == 1 && isNew {
	// 		return fmt.Errorf("缺少必填字段 %s", po.Variable)
	// 	}
	// }

	// data := customer.Extension
	// for key, value := range data {
	// 	po, ok := mt[key]
	// 	if !ok {
	// 		return fmt.Errorf("该字段在模板中不存在 %s", key)
	// 	}

	// 	val := reflect.ValueOf(value)
	// 	switch po.Kind {
	// 	case "String":
	// 		if val.Kind() != reflect.String {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "String", po.Variable, value)
	// 		}
	// 		if po.Reg == "" {
	// 			break
	// 		}

	// 		if ok := regexp.MustCompile(po.Reg).MatchString(value.(string)); !ok {
	// 			return fmt.Errorf("正则匹配失败 %v %s:%v", po.Reg, po.Variable, value)
	// 		}
	// 	case "Float32", "Float64":
	// 		if val.Kind() != reflect.Float32 && val.Kind() != reflect.Float64 {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "Float", po.Variable, value)
	// 		}
	// 	case "DateTime":
	// 		if val.Kind() != reflect.String {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "DateTime", po.Variable, value)
	// 		}
	// 		if _, err := time.Parse("2006-01-02 15:04:05", fmt.Sprint(value)); err != nil {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "DateTime", po.Variable, value)
	// 		}

	// 	case "Date":
	// 		if val.Kind() != reflect.String {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "Date", po.Variable, value)
	// 		}
	// 		if _, err := time.Parse("2006-01-02", fmt.Sprint(value)); err != nil {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "Date", po.Variable, value)
	// 		}
	// 	default:
	// 		ok := utils.InSlice([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
	// 			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
	// 			reflect.Uint64, reflect.Float32, reflect.Float64}, val.Kind())
	// 		if !ok {
	// 			return fmt.Errorf("错误类型 %v %s:%v", "Number", po.Variable, value)
	// 		}
	// 	}
	// }

	return nil
}
