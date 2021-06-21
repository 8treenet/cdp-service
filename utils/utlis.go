package utils

import (
	"strings"

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
