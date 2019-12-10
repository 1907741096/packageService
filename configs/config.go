package configs

import (
	lib "AS/commonlib-go"
	"AS/utils"
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path/filepath"
)



func initConfig(path string) {
	lib.InitConfig(path)
	initDataDB()
	//initLogDB()
	//initRedis()
	//initSentry()
	//initDomain()
	//initLog()
	initSystemSource()
	hooks := getLogHooks()
	//models.InitDB()
	lib.InitLogger(hooks)
}

func getLogHooks() []logrus.Hook {
	logPath := filepath.Join("/tmp")
	writer, err := rotatelogs.New(
		filepath.Join(logPath, "app.log"),
	)
	if err != nil {
		fmt.Errorf("config ratatelog fatal error. %v", errors.WithStack(err))
		return nil
	}

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.TraceLevel: writer,
	}
	formatter := &logrus.JSONFormatter{
		TimestampFormat: lib.DatetimeFormat,
	}

	hooks := []logrus.Hook{lfshook.NewHook(writeMap, formatter)}
	return hooks
}

func initDataDB() {
	utils.ParseJsonToStruct(&lib.DBConfig, lib.JsonData["database"].(map[string]interface{}))
	lib.DBConfig.URL = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		lib.DBConfig.User, lib.DBConfig.Password, lib.DBConfig.Host, lib.DBConfig.Port, lib.DBConfig.Database, lib.DBConfig.Charset)
	}

func initLogDB() {
	utils.ParseJsonToStruct(&lib.LogDBConfig, lib.JsonData["logdb"].(map[string]interface{}))
	lib.LogDBConfig.URL = fmt.Sprintf("%s:%d", lib.LogDBConfig.Host, lib.LogDBConfig.Port)
}

func initRedis() {
	utils.ParseJsonToStruct(&lib.RdsConfig, lib.JsonData["redis"].(map[string]interface{}))
	lib.RdsConfig.URL = fmt.Sprintf("%s:%d", lib.RdsConfig.Host, lib.RdsConfig.Port)
}

func initSentry() {
	utils.ParseJsonToStruct(&lib.SentryConfig, lib.JsonData["sentry"].(map[string]interface{}))
	lib.SentryConfig.URL = fmt.Sprintf("http://%s:%s@%s/%d",
		lib.SentryConfig.Key, lib.SentryConfig.Secret, lib.SentryConfig.Host, lib.SentryConfig.Project)
}

func initDomain() {
	lib.DomainConfig.Domain = lib.JsonData["domain"].(map[string]interface{})
}

func initLog() {
	utils.ParseJsonToStruct(&lib.LogConfig, lib.JsonData["log"].(map[string]interface{}))
}

func initSystemSource() {
	lib.SystemSource = lib.JsonData["systemSource"].(string)
}

//func initSignature() {
//	data := lib.JsonData["signature"].(map[string]interface{})
//	listData := map[string]services.KeyData{}
//	for key, value := range data{
//		var keyData services.KeyData
//		utils.ParseJsonToStruct(&keyData, value.(map[string]interface{}))
//		listData[key] = keyData
//	}
//	services.ConstSignatureConfig.C = listData
//}
