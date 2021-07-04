package cattle

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

	submit.metadata[ColumnCity] = ColumnTypeString
	submit.metadata[ColumnRegion] = ColumnTypeString
	submit.metadata[ColumnSourceId] = ColumnTypeInt16
	submit.metadata[ColumnIP] = ColumnTypeIP
	submit.metadata[ColumnCreateTime] = ColumnTypeDateTime
	submit.metadata[ColumnUserId] = ColumnTypeString
}

func (submit *Submit) SetLogger(l Logger) *Submit {
	submit.logger = l
	return submit
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
	case ColumnTypeString:
		if !ok {
			return "", nil
		}
		return datav, nil
	case ColumnTypeArrayString:
		if !ok {
			return clickhouse.Array([]string{}), nil
		}
		return datav, nil

	case ColumnTypeFloat32, ColumnTypeFloat64:
		if !ok {
			return 0.0, nil
		}
		return datav, nil
	case ColumnTypeArrayFloat32, ColumnTypeArrayFloat64:
		if !ok {
			return []float32{}, nil
		}
		return datav, nil
	case ColumnTypeUInt8, ColumnTypeUInt16, ColumnTypeUInt32, ColumnTypeUInt64, ColumnTypeInt8, ColumnTypeInt16, ColumnTypeInt32, ColumnTypeInt64:
		if !ok {
			return 0, nil
		}
		i, err := strconv.Atoi(fmt.Sprint(datav))
		if err != nil {
			return 0, err
		}
		return i, nil

	case ColumnTypeArrayUInt8, ColumnTypeArrayUInt16, ColumnTypeArrayUInt32, ColumnTypeArrayUInt64, ColumnTypeArrayInt8, ColumnTypeArrayInt16, ColumnTypeArrayInt32, ColumnTypeArrayInt64:
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

	case ColumnTypeDate:
		if !ok {
			return time.Now(), nil
		}
		return datav, nil
	case ColumnTypeArrayDate:
		if !ok {
			return []time.Time{time.Now()}, nil
		}
		return datav, nil
	case ColumnTypeDateTime:
		if !ok {
			return time.Now(), nil
		}
		return datav, nil
	case ColumnTypeArrayDateTime:
		if !ok {
			return []time.Time{time.Now()}, nil
		}
		return datav, nil
	case ColumnTypeIP:
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
