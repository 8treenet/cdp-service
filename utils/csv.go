package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
)

type CSV struct {
	header []string
	rows   [][]string
}

func NewCSV(data []byte) (*CSV, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	list, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("csv data empty.")
	}

	result := &CSV{}
	result.header = list[0]
	result.rows = list[1:]
	return result, nil
}

func (csv *CSV) VisitMap(f func(map[string]string)) {
	for _, row := range csv.rows {
		callValue := map[string]string{}
		for i := 0; i < len(csv.header); i++ {
			callValue[csv.header[i]] = row[i]
		}
		f(callValue)
	}
}

func (csv *CSV) ToMaps() (result []map[string]string) {
	result = make([]map[string]string, 0)
	for _, row := range csv.rows {
		m := map[string]string{}
		for i := 0; i < len(csv.header); i++ {
			m[csv.header[i]] = row[i]
		}
		result = append(result, m)
	}
	return
}

func (csv *CSV) VisitSlice(f func([]string)) {
	for _, row := range csv.rows {
		f(row)
	}
}
