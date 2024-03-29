package vo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type responseVO struct {
	code    int
	message string
	data    interface{}
}

func responseData(httpCode int, code int, message string, data interface{}, c *gin.Context) {
	c.JSON(httpCode, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func ResponseDataSuccess(data interface{}, c *gin.Context) {
	responseData(http.StatusOK, 0, "success", data, c)
}

func ResponseDataFail(c *gin.Context, err error) {
	responseData(http.StatusInternalServerError, 1, err.Error(), nil, c)
}

func ResponseWithoutData(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
