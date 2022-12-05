package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
)

// CreatePostHandler 创建并发布帖子
func CreatePostHandler(c *gin.Context) {

	//1. 获取并校验参数（从前端传来，发布的帖子内容、帖子的分区等数据）
	p := new(models.Post)
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

	userID, err := GetCurrentUser(c) //登陆时，鉴权中间件将User的ID保存到上下文中
	if err != nil {

		ResponseError(c, CodeUserNotLogin)
		return
	}
	p.AuthorID = userID

	//2. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	//3. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 查看对应URL中对应ID的帖子
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取并处理参数（post id）
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil { //处理字符，进制，大小
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 逻辑处理（在数据库中找post）
	data, err := logic.GetPostById(id)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应（对应的post）
	ResponseSuccess(c, data)
}

// GetPostDetailHandler 获取帖子列表，查看所有帖子
func GetPostListHandler(c *gin.Context) {
	//获取URL中的分页参数
	page, size := GetPageQuery(c)

	//获取列表数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// GetOrderPostListHandler
func GetOrderPostListHandler(c *gin.Context) {
	//请求参数是URL中的query 因此请求参数的结构体tag 使用 form
	//初始化结构体时指定初始参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("GetOrderPostListHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//获取列表数据
	data, err := logic.GetPostListInOrder(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
