// +build darwin debug
// +build !test
// +build !release

package configs

import (
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.DebugMode)
	initConfig("./configs/debug.json")
}