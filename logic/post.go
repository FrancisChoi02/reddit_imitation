package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {

	//1. 生成postID
	p.ID = snowflake.GenID()
	//2. 将post保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}

	err = redis.CreatePost(p.ID)
	if err != nil {
		return err
	}
	//3. 返回响应
	return
}

// GetPostById 根据ID获取帖子 （帖子由帖子内容本身 & 用户信息 & 帖子分区信息 组成）
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	//查询并组合我们接口想要的数据

	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed",
			zap.Int64("id", pid),
			zap.Error(err))
		return
	}

	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("id", post.AuthorID),
			zap.Error(err))
		return
	}

	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("id", post.CommunityID),
			zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostById 查询所有帖子 （帖子由帖子内容本身 & 用户信息 & 帖子分区信息 组成）
func GetPostList(page int64, size int64) (data []*models.ApiPostDetail, err error) {

	//和GetPostById的操作逻辑相似，只不过操作对象是 *models.ApiPostDetail 切片,用for循环解决
	//查询并组合该接口所需的数据
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	data = make([]*models.ApiPostDetail, 0, len(postList)) //为结果数组分配空间

	for _, post := range postList {

		//某一个数据获取错误不需要返回，记录在日志即可，允许返回空数据
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("id", post.CommunityID),
				zap.Error(err))
			continue
		}
		dataDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, dataDetail)
	}

	return
}
