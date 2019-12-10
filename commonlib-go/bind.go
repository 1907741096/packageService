package lib

import (
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
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
		// shouldBindBodyWith 在context存储了 BodyBytesKey, 需要在后置的middleware当中使用该值
		if err := c.ShouldBindBodyWith(obj, binding.JSON); err != nil && err != io.EOF {
			FormatFail(c, err)
			Warning("[bind]bind value error", ErrorMessage{"error": err.Error()}, "bind")
			c.Abort()
			return
		}
		c.Set(BindFormKey, obj)
		c.Next()
	}
}
