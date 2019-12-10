package lib

// log 初始化时使用hooks, 因为hooks当中的个性化配置的差异，所以讲hooks的配置独立
//func getLogHooks() []logrus.Hook {
//	logPath := filepath.Join("/data/logs/go", "feature")
//	writer, err := rotatelogs.New(
//		logPath+".%Y%m%d%H%M",
//		rotatelogs.WithLinkName(logPath),          // 生成软链，指向最新日志文件
//		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
//		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
//	)
//	if err != nil {
//		fmt.Errorf("config ratatelog fatal error. %v", errors.WithStack(err))
//		return nil
//	}
//
//	pathMap := lfshook.PathMap{
//		logrus.InfoLevel:  filepath.Join(logPath, "info.log"),
//		logrus.WarnLevel:  filepath.Join(logPath, "warning.log"),
//		logrus.ErrorLevel: filepath.Join(logPath, "error.log"),
//		logrus.TraceLevel: filepath.Join(logPath, "trace.log"),
//	}
//
//	writeMap := lfshook.WriterMap{
//		logrus.InfoLevel:  writer,
//		logrus.WarnLevel:  writer,
//		logrus.ErrorLevel: writer,
//		logrus.TraceLevel: writer,
//	}
//	formatter := &logrus.JSONFormatter{
//		TimestampFormat: lib.DatetimeFormat,
//	}
//
//	hooks := []logrus.Hook{lfshook.NewHook(pathMap, formatter), lfshook.NewHook(writeMap, formatter)}
//	return hooks
//}

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type Body struct {
	Message interface{}
}

type LogReport map[string]interface{}

var Logger *logrus.Logger

func InitLogger(hooks []logrus.Hook) {
	Logger = NewLogger(hooks)
}

func NewLogger(hooks []logrus.Hook) *logrus.Logger {
	if Logger != nil {
		return Logger
	}
	Logger = logrus.New()
	Logger.Level = logrus.TraceLevel

	for _, hook := range hooks {
		Logger.AddHook(hook)
	}
	return Logger
}

const DatetimeFormat = "2006-01-02 15:04:05"

type ErrorMessage map[string]interface{}

const LogInfo = "info"
const LogWarn = "warning"
const LogError = "error"

func Trace(body ...interface{}) {
	message := formatMessage(body)
	entry := setField("trace", "trace")
	report(body, "trace", "trace", LogInfo)
	entry.Info(message)
}

func Info(overview string, body interface{}, category string) {
	message := formatMessage(body)
	entry := setField(overview, category)
	report(body, overview, category, LogInfo)
	entry.Info(message)
}

func Warning(overview string, body interface{}, category string) {
	message := formatMessage(body)
	entry := setField(overview, category)
	report(body, overview, category, LogWarn)
	entry.Warning(message)
}

func Error(overview string, body interface{}, category string) {
	message := formatMessage(body)
	entry := setField(overview, category)
	report(body, overview, category, LogError)
	entry.Error(message)
}

func formatMessage(body interface{}) string {
	jsonData, err := json.Marshal(body)
	if err != nil {
		logrus.Error("format log message error: ", err.Error())
		return ""
	}

	return string(jsonData)
}

func setField(overview string, category string) *logrus.Entry {
	entry := Logger.WithFields(logrus.Fields{
		"overview": overview,
		"category": category,
	})

	return entry
}

// report sentry
func report(body interface{}, overview, category, level string) {
	reportData := map[string]interface{}{
		"message":  body,
		"category": category,
		"overview": overview,
		"level": level,
	}

	report := formatMessage(reportData)
	switch level {
	case LogInfo:
		return
	case LogWarn:
		packet := raven.NewPacketWithExtra(overview, raven.Extra{"content": reportData})
		packet.Level = raven.Severity(level)
		raven.Capture(packet, map[string]string{})
	case LogError:
		packet := raven.NewPacketWithExtra(overview, raven.Extra{"content": reportData}, raven.NewException(errors.New(report), raven.NewStacktrace(3, 3, nil)))
		packet.Level = raven.Severity(level)
		raven.Capture(packet, map[string]string{})
	}
}
