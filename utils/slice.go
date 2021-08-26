package utils

import (
	"errors"
	"reflect"
)

//InSlice .
func InSlice(slice interface{}, item interface{}) bool {
	values := reflect.ValueOf(slice)
	if values.Kind() != reflect.Slice {
		return false
	}

	size := values.Len()
	for index := 0; index < size; index++ {
		if values.Index(index).Interface() == item {
			return true
		}
	}
	return false
}

// 切割1维数组，target:传入数组, result:返回二维数据,capacity:切割的容量
func SliceUp(target interface{}, result interface{}, capacity int) error {
	if capacity <= 0 {
		return nil
	}
	targetValue := reflect.ValueOf(target)
	resultValue := reflect.ValueOf(result)

	if targetValue.Kind() != reflect.Slice {
		return errors.New("target type error")
	}

	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return errors.New("result type error")
	}

	newValue := reflect.MakeSlice(resultValue.Elem().Type(), 0, 0)
	begin := 0
	for {
		j := begin + capacity
		if j > targetValue.Len() {
			j = targetValue.Len()
		}

		rangeSlice := targetValue.Slice(begin, j)
		newValue = reflect.Append(newValue, rangeSlice)
		begin = j
		if j == targetValue.Len() {
			break
		}
	}

	resultValue.Elem().Set(newValue)
	return nil
}

func ToInterfaces(value interface{}) ([]interface{}, error) {
	result := []interface{}{}
	rvalue := reflect.ValueOf(value)
	if rvalue.Kind() != reflect.Slice {
		return result, errors.New("not slice")
	}
	for index := 0; index < rvalue.Len(); index++ {
		result = append(result, rvalue.Index(index).Interface())
	}
	return result, nil
}

//NewSlice .
func NewSlice(dsc interface{}, len int) error {
	dstv := reflect.ValueOf(dsc)
	if dstv.Elem().Kind() != reflect.Slice {
		return errors.New("dsc error")
	}

	result := reflect.MakeSlice(reflect.TypeOf(dsc).Elem(), len, len)
	dstv.Elem().Set(result)
	return nil
}
