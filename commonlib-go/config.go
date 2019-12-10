package lib

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var (
	JsonData     map[string]interface{}
	DBConfig     DatabaseConfig
	LogDBConfig  DatabaseConfig // LogDBConfig logdb 相关配置
	RdsConfig    RedisConfig    // RdsConfig redis 相关配置
	SentryConfig sentryConfig   // SentryConfig sentry 相关配置
	DomainConfig domainConfig   // 请求域名相关配置
	LogConfig    logConfig      // log相关配置
	SystemSource systemSource
)

// 解析配置文件
func InitConfig(path string) {
	config, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("read json file failed", err)
		os.Exit(-1)
	}

	stringConfig := regexp.MustCompile(`/\*.*\*/`).ReplaceAllString(string(config), "")
	if err := json.Unmarshal([]byte(stringConfig), &JsonData); err != nil {
		log.Println("invalid json in reading file", err)
		os.Exit(-1)
	}
}

// 数据库连接配置
type DatabaseConfig struct {
	Dialect        string
	Database       string
	User           string
	Password       string
	Host           string
	Port           int
	Charset        string
	URL            string
	MaxIdleConnNum int
	MaxOpenConnNum int
}


// sentry 配置
type sentryConfig struct {
	Host    string
	Key     string
	Secret  string
	Project int
	URL     string
}

// 域名配置
type domainConfig struct {
	Domain	map[string]interface{}
}

// 日志相关配置
type logConfig struct{
	Dir		 string
	Filename string
}

type systemSource = string

func RegisterDB(config DatabaseConfig) *gorm.DB{
	db, err := gorm.Open(config.Dialect, config.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if gin.Mode() == gin.ReleaseMode {
		db.SetLogger(logrus.StandardLogger())
	} else {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.MaxIdleConnNum)
	db.DB().SetMaxOpenConns(config.MaxOpenConnNum)
	db.InstantSet("gorm:association_autocreate", false)
	db.InstantSet("gorm:association_autoupdate", false)
	db.InstantSet("gorm:association_save_reference", false)
	return db
}