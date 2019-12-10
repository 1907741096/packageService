package main

import (
	_ "AS/configs"
	lib "AS/commonlib-go"
	"AS/controllers"
	"AS/models"
	"AS/middlewares"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	lib.RedisConnect(lib.RdsConfig)
	defer lib.RedisClose()

	router := getRouter()
	gracehttp.SetLogger(log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile))

	gracehttp.Serve(
		&http.Server{Addr: ":80", Handler: router},
	)
}

func getRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), lib.MonitorFunc())

	// 前端
	v1 := router.Group("/v1")
	{
		v1.POST("/v/s", middlewares.Bind(models.GetOneDataForm{}), controllers.GetStatus)
	}

	// 后台
	internal := router.Group("/internal")
	{
		internal.POST("/app/get-version-list", middlewares.Bind(models.GetVersionListDataForm{}), controllers.GetVersionList)
		internal.POST("/app/get-package-list", middlewares.Bind(models.GetPackageListDataForm{}), controllers.GetPackageList)
		internal.POST("/app/save-version", middlewares.Bind(models.SaveVersionDataForm{}), controllers.SaveVersion)
		internal.POST("/app/save-package", middlewares.Bind(models.SavePackageDataForm{}), controllers.SavePackage)
	}

	router.HEAD("/ping", func(context *gin.Context) {
		lib.FormatSuccess(context, gin.H{})
	})
	return router
}