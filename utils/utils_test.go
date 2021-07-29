package utils

import (
	"fmt"
	"net"
	"strconv"
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
	t.Log(strconv.ParseFloat("1", 64))
}
