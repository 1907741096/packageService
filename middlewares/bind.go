package middlewares

import (
	"AS/utils"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

func Bind(val interface{}) gin.HandlerFunc {
	value := reflect.ValueOf(val)
	if value.Kind() == reflect.Ptr {
		panic(`Bind struct can not be a pointer. Example:
	Use: gin.Bind(Struct{}) instead of gin.Bind(&Struct{})
`)
	}

	return func(c *gin.Context) {
		obj := reflect.New(value.Type()).Interface()
		if err := c.ShouldBind(obj); err != nil {
			log.Println(err)
			utils.FormatFail(c, utils.ErrParamInvalid)
			c.Abort()
			return
		}
		c.Set(utils.BindFormKey, obj)
		c.Next()
	}
}
