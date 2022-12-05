package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

//本来在middleware包 避免循环引用
const CtxUserID = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// 获取当前登录的用户IP
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	//从上下文中获取用户ID
	//用户登陆时，中间件获取用户ID并保存到CtxUserID变量中
	uid, ok := c.Get(CtxUserID)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

func GetPageQuery(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 0
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
