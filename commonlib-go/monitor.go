package lib

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func MonitorFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		request, ok := c.Get(gin.BodyBytesKey)
		var requestBody []byte
		if ok {
			if requestBytes, ok := request.([]byte); ok {
				requestBody = requestBytes
			}
		}

		response, ok := c.Get(ResponseDataKey)
		responseData, typeOf := response.(gin.H)

		if !ok || !typeOf {
			Warning("[monitor]processing monitor middleware error", ErrorMessage{
				"error": ErrTypeTransErr,
			}, "monitor")
			return
		}

		requestString := strings.Replace(strings.Replace(string(requestBody), "\n", "", -1), "\\", "", -1)
		Info("[monitor]http request message", ErrorMessage{
			"request": strings.Replace(requestString, "\t", "", -1),
			"response": responseData,
		}, "monitor")
	}
}
