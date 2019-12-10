package utils

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

const (
	BindFormKey     = "_controllers/bind_form_key"
	ResponseDataKey = "_controllers/response_data_key"
)

func FormatSuccess(c *gin.Context, data interface{}) {
	data = gin.H{
		"code":    "0",
		"message": "success",
		"data":    data,
	}
	c.Set(ResponseDataKey, data)
	c.JSON(http.StatusOK, data)
}


func FormatFail(c *gin.Context, err error) {
	code, ok := ErrCodeMap[err]
	if !ok {
		code = SystemInterError
	}

	data := gin.H{
		"code":    code,
		"message": err.Error(),
		"data":    gin.H{},
	}

	c.Set(ResponseDataKey, data)
	c.JSON(http.StatusOK, data)
}