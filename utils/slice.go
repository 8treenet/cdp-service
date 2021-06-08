package utils

import "reflect"

//InSlice .
func InSlice(array interface{}, item interface{}) bool {
	values := reflect.ValueOf(array)
	if values.Kind() != reflect.Slice {
		return false
	}

	size := values.Len()
	list := make([]interface{}, size)
	slice := values.Slice(0, size)
	for index := 0; index < size; index++ {
		list[index] = slice.Index(index).Interface()
	}

	for index := 0; index < len(list); index++ {
		if list[index] == item {
			return true
		}
	}
	return false
}
