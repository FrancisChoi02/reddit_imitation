package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//用于 给前端返回消息 的结构体
type ResponseData struct {
	Code ResCode     `json:"code"` //程序的错误码
	Msg  interface{} `json:"msg"`  //提示信息
	Data interface{} `json:"data"` //数据
}

// ResponseAuthenticationError 返回请求头格式错误
func ResponseAuthenticationError(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 2004,
		"msg":  "请求头中auth格式有误",
	})
}

// ResponseError 返回错误消息
func ResponseError(c *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(), //调用ResCode的方法  获取状态码对应的提示字串
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseErrorWithMsg 返回带自定义提示的错误消息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg, //不适用code对应的提示字串，而是用传入的
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

// ResponseSuccess 返回正确消息
func ResponseSuccess(c *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(), //调用ResCode的方法  获取状态码对应的提示字串
		Data: data,
	}
	c.JSON(http.StatusOK, rd)
}
