package utils

import (
	"bytes"
	"runtime"
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
