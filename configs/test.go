// +build test

package configs

import (
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	initConfig("./configs/test.json")
}
