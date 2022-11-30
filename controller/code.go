package controller

type ResCode int64

//自增状态码常量
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeUserNotLogin
	CodeAuthError
	CodeTokenInvalid
)

//状态码 对应的提示语句
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeUserNotLogin: "请登录",
	CodeAuthError:    "请求头中auth格式有误",
	CodeTokenInvalid: "Token无效",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy] //非定义好的范围内的状态码，一律返回服务繁忙
	}
	return msg
}
