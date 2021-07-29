package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	uuid "github.com/iris-contrib/go.uuid"
)

// GenerateUUID Build a unique ID.
func GenerateUUID() (string, error) {
	uuidv1, e := uuid.NewV1()
	if e != nil {
		return "", e
	}
	return strings.ReplaceAll(uuidv1.String(), "-", ""), nil
}

func DateTimeFormat(t time.Time) string {
	return t.Local().Format("2006-01-02 15:04:05")
}

func GetStackTrace() []byte {
	def := []byte("An unknown error")
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, 2<<10) //2KB
	length := runtime.Stack(stack, false)
	start := bytes.Index(stack, s)
	if start == -1 || start > len(stack) || length > len(stack) {
		return def
	}
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	if start >= len(stack) {
		return def
	}
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end > len(stack) {
		return def
	}
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end > len(stack) {
		return def
	}
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return stack
}

func IsNumber(num interface{}) bool {
	strNum := fmt.Sprint(num)
	_, intErr := strconv.ParseInt(strNum, 10, 64)
	_, uintErr := strconv.ParseUint(strNum, 10, 64)
	_, floatErr := strconv.ParseFloat(strNum, 64)
	if intErr != nil && floatErr != nil && uintErr != nil {
		return false
	}
	return true
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

func IsNumberSlice(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < v.Len(); i++ {
		if !IsNumber(v.Index(i)) {
			return false
		}
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
