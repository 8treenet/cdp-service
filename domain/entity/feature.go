package entity

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"cdp-service/domain/po"
	"cdp-service/infra/cattle"
	"cdp-service/utils"

	"github.com/8treenet/freedom"
)

// Feature 特征实体
type Feature struct {
	po.BehaviourFeature
	freedom.Entity
	FeatureMetadata []*po.BehaviourFeatureMetadata
}

func (entity *Feature) Identity() string {
	return fmt.Sprint(entity.ID)
}

func (entity *Feature) AddMetadata(variable, title, kind, dict string, orderByNumber int) {
	entity.FeatureMetadata = append(entity.FeatureMetadata, &po.BehaviourFeatureMetadata{
		Variable:      variable,
		Title:         title,
		Kind:          kind,
		Dict:          dict,
		OrderByNumber: orderByNumber,
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

func (entity *Feature) GetWarehouse() string {
	return entity.Warehouse
}

func (entity *Feature) GetColumnType(column string) string {
	for _, v := range entity.FeatureMetadata {
		if v.Variable == column {
			return v.Kind
		}
	}
	return ""
}

func (entity *Feature) CheckMetadata(data map[string]interface{}) error {
	for _, v := range entity.FeatureMetadata {
		value, ok := data[v.Variable]
		if !ok {
			continue
		}
		typeErr := fmt.Errorf("feature:%v variable %s value:%v, unable.", v.Title, v.Variable, value)

		rv := reflect.ValueOf(value)
		switch v.Kind {
		case cattle.ColumnTypeDate, cattle.ColumnTypeDateTime:
			if !utils.IsDateTime(fmt.Sprint(value)) {
				return typeErr
			}
		case cattle.ColumnTypeString:
			if rv.Kind() != reflect.String {
				return typeErr
			}
		case cattle.ColumnTypeFloat32, cattle.ColumnTypeFloat64:
			if _, e := utils.ToFloat(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeUInt8, cattle.ColumnTypeUInt16, cattle.ColumnTypeUInt32, cattle.ColumnTypeUInt64:
			if _, e := utils.ToUint(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeInt8, cattle.ColumnTypeInt16, cattle.ColumnTypeInt32, cattle.ColumnTypeInt64:
			if _, e := utils.ToInt(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeArrayFloat32, cattle.ColumnTypeArrayFloat64:
			if _, e := utils.ToFloatSlice(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeArrayUInt8, cattle.ColumnTypeArrayUInt16, cattle.ColumnTypeArrayUInt32, cattle.ColumnTypeArrayUInt64:
			if _, e := utils.ToUintSlice(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeArrayInt8, cattle.ColumnTypeArrayInt16, cattle.ColumnTypeArrayInt32, cattle.ColumnTypeArrayInt64:
			if _, e := utils.ToIntSlice(value); e != nil {
				return typeErr
			}
		case cattle.ColumnTypeArrayString:
			if rv.Kind() != reflect.Slice {
				return typeErr
			}
		case cattle.ColumnTypeArrayDate, cattle.ColumnTypeArrayDateTime:
			if !utils.IsDateTimeSlice(value) {
				return typeErr
			}
		case cattle.ColumnTypeIP:
			if net.ParseIP(fmt.Sprint(value)) == nil {
				return typeErr
			}
		}
	}

	return nil
}

func (entity *Feature) ConvertMetadata(data map[string]string) (result map[string]interface{}, e error) {
	result = make(map[string]interface{})

	for _, v := range entity.FeatureMetadata {
		value, ok := data[v.Variable]
		if !ok {
			continue
		}
		typeErr := fmt.Errorf("feature:%v variable %s value:%v, unable.", v.Title, v.Variable, value)

		switch v.Kind {
		case cattle.ColumnTypeDate, cattle.ColumnTypeDateTime:
			if !utils.IsDateTime(fmt.Sprint(value)) {
				e = typeErr
				return
			}
			result[v.Variable] = value

		case cattle.ColumnTypeFloat32, cattle.ColumnTypeFloat64:
			fnum, err := utils.ToFloat(value)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = fnum
		case cattle.ColumnTypeUInt8, cattle.ColumnTypeUInt16, cattle.ColumnTypeUInt32, cattle.ColumnTypeUInt64:
			unum, err := utils.ToUint(value)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = unum
		case cattle.ColumnTypeInt8, cattle.ColumnTypeInt16, cattle.ColumnTypeInt32, cattle.ColumnTypeInt64:
			unum, err := utils.ToInt(value)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = unum

		case cattle.ColumnTypeArrayFloat32, cattle.ColumnTypeArrayFloat64:
			list := strings.Split(value, "|")
			flist, err := utils.ToFloatSlice(list)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = flist

		case cattle.ColumnTypeArrayUInt8, cattle.ColumnTypeArrayUInt16, cattle.ColumnTypeArrayUInt32, cattle.ColumnTypeArrayUInt64:
			list := strings.Split(value, "|")
			flist, err := utils.ToUintSlice(list)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = flist
		case cattle.ColumnTypeArrayInt8, cattle.ColumnTypeArrayInt16, cattle.ColumnTypeArrayInt32, cattle.ColumnTypeArrayInt64:
			list := strings.Split(value, "|")
			flist, err := utils.ToIntSlice(list)
			if err != nil {
				e = typeErr
				return
			}
			result[v.Variable] = flist

		case cattle.ColumnTypeArrayString:
			list := strings.Split(value, "|")
			result[v.Variable] = list
		case cattle.ColumnTypeArrayDate, cattle.ColumnTypeArrayDateTime:
			list := strings.Split(value, "|")
			if !utils.IsDateTimeSlice(list) {
				e = typeErr
				return
			}
			result[v.Variable] = list
		case cattle.ColumnTypeIP:
			if net.ParseIP(value) == nil {
				e = typeErr
				return
			}
			result[v.Variable] = value
		default:
			result[v.Variable] = value
		}
	}

	return
}
