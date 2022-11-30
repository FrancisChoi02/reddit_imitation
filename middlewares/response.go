package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseAuthenticationError 返回请求头格式错误
func ResponseAuthenticationError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 2004,
		"msg":  "请求头中auth格式有误",
	})
}

func ResponseTokenError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 2005,
		"msg":  "无效的Token",
	})
}
