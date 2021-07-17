package cattle

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ColumnRegion     = "region"
	ColumnCity       = "city"
	ColumnIP         = "ip"
	ColumnSourceId   = "sourceId"
	ColumnCreateTime = "createTime"
	ColumnUserId     = "userId"
)

const (
	ColumnTypeString      = "String"
	ColumnTypeArrayString = "ArrayString"

	ColumnTypeFloat32 = "Float32"
	ColumnTypeFloat64 = "Float64"

	ColumnTypeArrayFloat32 = "ArrayFloat32"
	ColumnTypeArrayFloat64 = "ArrayFloat64"

	ColumnTypeUInt8  = "UInt8"
	ColumnTypeUInt16 = "UInt16"
	ColumnTypeUInt32 = "UInt32"
	ColumnTypeUInt64 = "UInt64"

	ColumnTypeArrayUInt8  = "ArrayUInt8"
	ColumnTypeArrayUInt16 = "ArrayUInt16"
	ColumnTypeArrayUInt32 = "ArrayUInt32"
	ColumnTypeArrayUInt64 = "ArrayUInt64"

	ColumnTypeInt8  = "Int8"
	ColumnTypeInt16 = "Int16"
	ColumnTypeInt32 = "Int32"
	ColumnTypeInt64 = "Int64"

	ColumnTypeArrayInt8  = "ArrayInt8"
	ColumnTypeArrayInt16 = "ArrayInt16"
	ColumnTypeArrayInt32 = "ArrayInt32"
	ColumnTypeArrayInt64 = "ArrayInt64"

	ColumnTypeDateTime = "DateTime"
	ColumnTypeDate     = "Date"

	ColumnTypeArrayDateTime = "ArrayDateTime"
	ColumnTypeArrayDate     = "ArrayDate"

	ColumnTypeIP = "IPv4"
)

func ArrayKind(kind string) string {
	list := strings.Split(kind, "Array")
	if len(list) == 2 {
		return fmt.Sprintf("%s(%s)", "Array", list[1])
	}
	return kind
}

func toNumber(value interface{}) (result interface{}) {
	result = value
	str := fmt.Sprint(value)

	if strings.Contains(str, ".") {
		fvalue, err := strconv.ParseFloat(str, 64)
		if err == nil {
			result = fvalue
			return
		}
	}

	if iv, err := strconv.ParseInt(str, 10, 64); err == nil {
		result = iv
		return
	}

	if uiv, err := strconv.ParseUint(str, 10, 64); err == nil {
		result = uiv
		return
	}
	return
}
