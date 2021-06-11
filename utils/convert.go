package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ToInt(value interface{}) (int, error) {
	str := fmt.Sprint(value)
	list := strings.Split(str, ".")
	if len(list) == 0 {
		return 0, errors.New("not int")
	}
	return strconv.Atoi(list[0])
}
