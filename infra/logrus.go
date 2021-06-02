package infra

import (
	"time"

	"github.com/8treenet/freedom"
	"github.com/kataras/golog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func NewLogrusMiddleware(logFolder string, console bool) func(value *freedom.LogRow) bool {
	initLogger(logFolder)
	return func(value *freedom.LogRow) bool {
		level := toLogrusLevel(value.Level)
		loggerEntity.WithFields(logrus.Fields(value.Fields)).Log(level, value.Message)
		return !console // 返回true 停止中间件遍历，最底层默认console
	}
}

var loggerEntity *logrus.Logger

func initLogger(logFolder string) {
	loggerEntity = logrus.New()
	loggerEntity.SetLevel(logrus.DebugLevel)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
	loggerEntity.SetFormatter(customFormatter)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		logFolder+"/%Y%m%d.log",
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
