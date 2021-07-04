package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/8treenet/freedom"
	"github.com/ClickHouse/clickhouse-go"
)

type Submit struct {
	logger    Logger
	manager   *Manager
	tableName string
	metadata  map[string]string
	rows      []map[string]interface{}
	defValue  struct {
		region   string
		city     string
		ip       string
		sourceId int
	}
}

func (submit *Submit) init() {
	submit.logger = freedom.Logger()
	submit.metadata = make(map[string]string)
	submit.rows = make([]map[string]interface{}, 0)

	submit.metadata["city"] = "String"
	submit.metadata["region"] = "String"
	submit.metadata["sourceId"] = "Int16"
	submit.metadata["ip"] = "IPv4"
	submit.metadata["createTime"] = "DateTime"
}

func (submit *Submit) SetLogger(l Logger) {
	submit.logger = l
}

func (submit *Submit) Do() error {
	keys := []string{}
	perchs := []string{}
	for variable := range submit.metadata {
		keys = append(keys, variable)
		perchs = append(perchs, "?")
	}
	parpare := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", submit.tableName, strings.Join(keys, ","), strings.Join(perchs, ","))

	return submit.manager.tx(func(tx *sql.Tx) error {
		stmt, ferr := tx.Prepare(parpare)
		if ferr != nil {
			return ferr
		}
		return submit.stmtAdd(stmt, keys)
	})
}

func (submit *Submit) stmtAdd(stmt *sql.Stmt, keys []string) error {
	for _, data := range submit.rows {
		var args []interface{}
		for _, columnName := range keys {
			one, e := submit.parse(data, columnName)
			if e != nil {
				submit.logger.Errorf("stmtAdd data:%v, columnName:%v, error:%v", data, columnName, e)
				continue
			}
			args = append(args, one)
		}

		if _, err := stmt.Exec(args...); err != nil {
			return err
		}
	}
	return nil
}

func (submit *Submit) parse(data map[string]interface{}, columnName string) (interface{}, error) {
	datav, ok := data[columnName]
	kind := submit.metadata[columnName]
	switch kind {
	case "String":
		if !ok {
			return "", nil
		}
		return datav, nil
	case "ArrayString":
		if !ok {
			return clickhouse.Array([]string{}), nil
		}
		return datav, nil

	case "Float32", "Float64":
		if !ok {
			return 0.0, nil
		}
		return datav, nil
	case "ArrayFloat32", "ArrayFloat64":
		if !ok {
			return []float32{}, nil
		}
		return datav, nil
	case "UInt8", "UInt16", "UInt32", "UInt64", "Int8", "Int16", "Int32", "Int64":
		if !ok {
			return 0, nil
		}
		i, err := strconv.Atoi(fmt.Sprint(datav))
		if err != nil {
			return 0, err
		}
		return i, nil

	case "ArrayUInt8", "ArrayUInt16", "ArrayUInt32", "ArrayUInt64", "ArrayInt8", "ArrayInt16", "ArrayInt32", "ArrayInt64":
		if !ok {
			return clickhouse.Array([]int{}), nil
		}
		listValue := reflect.ValueOf(datav)
		if listValue.Kind() != reflect.Slice {
			return nil, errors.New("该数据不是数组")
		}

		newDatav := []int{}
		for i := 0; i < listValue.Len(); i++ {
			numstr := strings.Split(fmt.Sprint(listValue.Index(i).Interface()), ".")[0]
			i, _ := strconv.Atoi(numstr)
			newDatav = append(newDatav, i)
		}
		return newDatav, nil

	case "Date":
		if !ok {
			return time.Now(), nil
		}
		return time.ParseInLocation("2006-01-02", fmt.Sprint(datav), time.Local)
	case "ArrayDate":
		if !ok {
			return []time.Time{time.Now()}, nil
		}
		// datetimes, cok := datav.([]string)
		// if !cok {
		// 	return nil, fmt.Errorf("ArrayDate error")
		// }
		// var newDatav []time.Time
		// for _, v := range datetimes {
		// 	datev, err := time.ParseInLocation("2006-01-02", fmt.Sprint(v), time.Local)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	newDatav = append(newDatav, datev)
		// }
		//return newDatav, nil
		return datav, nil
	case "DateTime":
		if !ok {
			return time.Now(), nil
		}
		return time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(datav), time.Local)
	case "ArrayDateTime":
		if !ok {
			return []time.Time{time.Now()}, nil
		}
		// datetimes, cok := datav.([]string)
		// if !cok {
		// 	return nil, fmt.Errorf("ArrayDateTime error")
		// }
		// var newDatav []time.Time
		// for _, v := range datetimes {
		// 	datetimev, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprint(v), time.Local)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	newDatav = append(newDatav, datetimev)
		// }
		// return newDatav, nil
		return datav, nil
	case "IPv4":
		if !ok {
			return "0.0.0.0", nil
		}
		return datav, nil
	}

	return nil, errors.New("未知类型")
}

func (submit *Submit) AddRow(row map[string]interface{}) {
	submit.rows = append(submit.rows, row)
}

func (submit *Submit) AddMetadata(variable, kind string) {
	submit.metadata[variable] = kind
}
