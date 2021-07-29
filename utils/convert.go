package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToFloat(num interface{}) (float64, error) {
	strNum := fmt.Sprint(num)
	return strconv.ParseFloat(strNum, 64)
}

func ToFloatSlice(array interface{}) ([]float64, error) {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("not slice")
	}
	result := []float64{}
	for i := 0; i < v.Len(); i++ {
		str := fmt.Sprint(v.Index(i).Interface())
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func ToInt(num interface{}) (int64, error) {
	strNum := fmt.Sprint(num)
	return strconv.ParseInt(strNum, 10, 64)
}

func ToIntSlice(array interface{}) ([]int64, error) {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("not slice")
	}
	result := []int64{}
	for i := 0; i < v.Len(); i++ {
		str := fmt.Sprint(v.Index(i).Interface())
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func ToUint(num interface{}) (uint64, error) {
	strNum := fmt.Sprint(num)
	return strconv.ParseUint(strNum, 10, 64)
}

func ToUintSlice(array interface{}) ([]uint64, error) {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("not slice")
	}
	result := []uint64{}
	for i := 0; i < v.Len(); i++ {
		str := fmt.Sprint(v.Index(i).Interface())
		num, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

func ToNumber(strNum string) interface{} {
	i, intErr := strconv.ParseInt(strNum, 10, 64)
	if intErr == nil {
		return i
	}
	ui, uintErr := strconv.ParseUint(strNum, 10, 64)
	if uintErr == nil {
		return ui
	}
	f, floatErr := strconv.ParseFloat(strNum, 64)
	if floatErr == nil {
		return f
	}
	return 0
}

func IsDateTime(value string) bool {
	_, dtErr := time.Parse("2006-01-02 15:04:05", value)
	_, dErr := time.Parse("2006-01-02", value)
	if dtErr != nil && dErr != nil {
		return false
	}
	return true
}

func IsDateTimeSlice(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		if !IsDateTime(fmt.Sprint(v.Index(i))) {
			return false
		}
	}
	return true
}

func ToStringSlice(array interface{}) ([]string, error) {
	v := reflect.ValueOf(array)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("not slice")
	}
	result := []string{}
	for i := 0; i < v.Len(); i++ {
		str := fmt.Sprint(v.Index(i).Interface())
		result = append(result, str)
	}
	return result, nil
}
