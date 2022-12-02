package controller

import (
	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

//与社区相关的业务分发

// CommunityHandler 查询所有社区
func CommunityHandler(c *gin.Context) {
	//查询所有的社区，以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //服务端的报错不可以暴露给用户
		return
	}
	//如果社区列表获取成功，将数据返回给前端
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 查询某个ID对应的社区详情
func CommunityDetailHandler(c *gin.Context) {
	//1.获取URL中的参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil { //处理字符，进制，大小
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.业务逻辑处理
	data, err := logic.GetCommunityDetail(id)

	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //服务端的报错不可以暴露给用户
		return
	}
	//3.结果返回
	ResponseSuccess(c, data)

}
