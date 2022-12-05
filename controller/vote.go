package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PostVoteController(c *gin.Context) {
	p := new(models.ParmVoteData)
	if err := c.ShouldBind(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //ParmVoteData中涉及投票的参数成员靠validator检验
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 翻译并去除错误提示中的 结构体标识
		errData := removeTopStruct(errs.Translate(trans)) //trans是全局的翻译器
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}

	// 获取当前请求的用户的id
	userID, err := GetCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeUserNotLogin)
		return
	}

	// 具体的投票业务逻辑
	logic.VoteForPost(userID, p)

	ResponseSuccess(c, nil)
}
