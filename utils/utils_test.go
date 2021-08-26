package utils

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestUtils(t *testing.T) {
	data := []byte(`1,2,3,4,5
2,3,4,5,6
7,8,9,10,9`)
	csv, err := NewCSV(data)
	if err != nil {
		panic(err)
	}
	csv.VisitMap(func(m map[string]string) {
		t.Log(m["1"])
	})
	t.Log(csv.ToMaps())

	i := net.ParseIP(fmt.Sprint("200.200.200.200"))
	t.Log(i)
	//_, err = time.Parse("2006-01-02 15:04:05", "2021-01-02")
	//t.Log(err)
}

func TestUtils2(t *testing.T) {
	vv := []string{"controller", "service", "repository"}
	r := reflect.ValueOf(vv)
	t.Log(r.Type().Elem().Kind())
}

func TestInSlice(t *testing.T) {
	if !InSlice([]string{"1", "2", "3", "4"}, "1") {
		panic("")
	}

	if !InSlice([]int{1, 2, 3, 4, 5}, 3) {
		panic("")
	}
	t.Log(InSlice([]int{1, 2, 3, 4, 5}, 3))
}
