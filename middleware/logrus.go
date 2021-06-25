package middleware

import (
	"fmt"
	"sort"
	"time"

	"github.com/8treenet/freedom"
	"github.com/kataras/golog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type logrusRow struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
}

func NewLogrusMiddleware(logFolder string, console bool) func(value *freedom.LogRow) bool {
	initLogger(logFolder)
	write := make(chan *logrusRow, 3000)
	go func() {
		for row := range write {
			//使用Logrus输出
			loggerEntity.WithFields(row.Fields).Log(row.Level, row.Message)
		}
	}()

	return func(value *freedom.LogRow) bool {
		//打入管道
		write <- &logrusRow{
			Level:   toLogrusLevel(value.Level),
			Message: value.Message,
			Fields:  logrus.Fields(value.Fields),
		}

		if !console {
			return true // 返回true 停止中间件遍历，最底层默认console
		}

		//组织console输出
		fieldKeys := []string{}
		for k := range value.Fields {
			fieldKeys = append(fieldKeys, k)
		}
		sort.Strings(fieldKeys)
		for i := 0; i < len(fieldKeys); i++ {
			fieldMsg := value.Fields[fieldKeys[i]]
			if value.Message != "" {
				value.Message += "  "
			}
			value.Message += fmt.Sprintf("%s=%v", fieldKeys[i], fieldMsg)
		}
		return false // 返回false 继续中间件遍历，最底层默认console
	}
}

var loggerEntity *logrus.Logger

func initLogger(logFolder string) {
	loggerEntity = logrus.New()
	loggerEntity.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.JSONFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
	loggerEntity.SetFormatter(customFormatter)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		logFolder+"/%Y_%m_%d.log",
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		freedom.Logger().Fatal(err)
	}
	loggerEntity.SetOutput(logWriter)
}

func toLogrusLevel(level golog.Level) logrus.Level {
	switch level {
	case golog.DebugLevel:
		return logrus.DebugLevel
	case golog.ErrorLevel:
		return logrus.ErrorLevel
	case golog.InfoLevel:
		return logrus.InfoLevel
	case golog.FatalLevel:
		return logrus.FatalLevel
	case golog.WarnLevel:
		return logrus.WarnLevel
	}
	return logrus.InfoLevel
}
