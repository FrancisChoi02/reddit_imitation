package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//1.获取参数 和 参数校验
	var p models.ParamSignUp //初始化用于获取请求参数的结构体
	if err := c.ShouldBind(&p); err != nil {

		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam) //状态码为 参数错误
			return
		}

		// validator.ValidationErrors类型错误则进行翻译
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//去除报告中 错误位置提示的前缀
		tmpString := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, tmpString)
		return

	}

	fmt.Println(p) //打印测试参数获取情况
	//手动对请求参数进行消息的业务规则校验

	//2.业务逻辑处理
	if err := logic.SignUp(&p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) { //mysql包里面的ErrorUserExist报错字串
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)

}

func LoginHandler(c *gin.Context) {
	//1.获取参数 和 参数校验
	var p models.ParamLogin //初始化用于获取请求参数的结构体
	if err := c.ShouldBind(&p); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam) //状态码为 参数错误
			return
		}

		// validator.ValidationErrors类型错误则进行翻译
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//去除报告中 错误位置提示的前缀
		tmpString := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, tmpString)
		return
	}
	//2.业务逻辑处理
	token, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("SignUp with invalid param", zap.String("username", p.Username), zap.Error(err))
		//返回登陆失败消息
		if errors.Is(err, mysql.ErrorUserNotExist) { //mysql包里面的ErrorUserExist报错字串
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, token) //将token字符串放到赋值给data成员变量传给前端
}
