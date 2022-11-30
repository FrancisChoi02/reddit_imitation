package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token的方式： 放在请求头
		// Token放在Header的Authorization中，并使用Bearer开头
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			controller.ResponseError(c, controller.CodeUserNotLogin)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			//返回请求头错误
			controller.ResponseError(c, controller.CodeAuthError)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		// mc是从jwt中解析出来的 装着用户信息的MyClaim结构体
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			//返回token无效
			controller.ResponseError(c, controller.CodeTokenInvalid)
			c.Abort()
			return
		}

		// 将当前请求的username信息 保存到请求的上下文c上
		c.Set(controller.CtxUserID, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get(CtxUserID)来获取当前请求的用户信息
	}
}
