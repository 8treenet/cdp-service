package cattle

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"cdp-service/utils"

	"github.com/go-xorm/builder"
)

type Warehouse interface {
	GetWarehouse() string
	GetColumnType(string) string
}

func (dsl *DSL) arrayIn(from, column string, listValue interface{}) builder.Cond {
	value := dsl.convertValue(from, column, listValue)
	if reflect.ValueOf(value).Kind() == reflect.Slice {
		return builder.Expr(fmt.Sprintf("hasAny(%s.%s,[%s])", from, column, listValue))
	}

	list, _ := utils.ToInterfaces(strings.Split(listValue.(string), ","))
	typeList := []interface{}{}
	for _, v := range list {
		typeList = append(typeList, dsl.convertValue(from, column, v))
	}

	return builder.In(from+"."+column, typeList...)
}

func (dsl *DSL) SetMetedata(w Warehouse) {
	dsl.whs = append(dsl.whs, w)
}

func (dsl *DSL) convertValue(from, column string, value interface{}) (result interface{}) {
	result = fmt.Sprint(value)
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Struct {
		return value
	}
	if column == ColumnSourceId {
		if ivalue, e := strconv.ParseInt(fmt.Sprint(value), 10, 64); e == nil {
			result = ivalue
			return
		}
	}

	for _, v := range dsl.whs {
		if v.GetWarehouse() != from {
			continue
		}

		switch v.GetColumnType(column) {
		case ColumnTypeArrayString:
			strList := strings.Split(fmt.Sprint(value), ",")
			result = strList
		case ColumnTypeFloat32, ColumnTypeFloat64:
			if fvalue, e := strconv.ParseFloat(fmt.Sprint(value), 64); e == nil {
				result = fvalue
				return
			}
		case ColumnTypeArrayFloat32, ColumnTypeArrayFloat64:
			strList := strings.Split(fmt.Sprint(value), ",")
			fList := []float64{}
			for _, v := range strList {
				fvalue, e := strconv.ParseFloat(v, 64)
				if e != nil {
					break
				}
				fList = append(fList, fvalue)
			}
			result = fList

		case ColumnTypeUInt8, ColumnTypeUInt16, ColumnTypeUInt32, ColumnTypeUInt64, ColumnTypeInt8, ColumnTypeInt16, ColumnTypeInt32, ColumnTypeInt64:
			if ivalue, e := strconv.ParseInt(fmt.Sprint(value), 10, 64); e == nil {
				result = ivalue
				return
			}
		case ColumnTypeArrayUInt8, ColumnTypeArrayUInt16, ColumnTypeArrayUInt32, ColumnTypeArrayUInt64, ColumnTypeArrayInt8, ColumnTypeArrayInt16, ColumnTypeArrayInt32, ColumnTypeArrayInt64:
			strList := strings.Split(fmt.Sprint(value), ",")
			intList := []int64{}
			for _, v := range strList {
				ivalue, e := strconv.ParseInt(v, 10, 64)
				if e != nil {
					break
				}
				intList = append(intList, ivalue)
			}
			result = intList
		}

		break
	}
	return result
}
