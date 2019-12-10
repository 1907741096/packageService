// +build release

package configs

import (
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	initConfig("/usr/local/services/app/configs/release.json")
}
