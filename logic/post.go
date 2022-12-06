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

	err = redis.CreatePost(p.ID, p.CommunityID)
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

// GetPostListInOrder 根据Order参数表达的方式  返回所有帖子
func GetPostListInOrder(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 1.根据p.Order从redis中获取对应的id列表（key）
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) returned 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids)) //将结果打印到日志里面，检查id列表获取情况

	// 2.根据id列表，从MySQL中查询帖子的详细信息
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids)) //将结果打印到日志里面，检查帖子列表的获取情况

	// 查询每篇帖子的 投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 与GetPostList()的处理方式相同，组装结果切片
	//将帖子作者、分区信息 查询出并填充到帖子结果切片中
	data = make([]*models.ApiPostDetail, 0, len(postList)) //为结果数组分配空间

	for i, post := range postList {

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
			VoteNum:         votedata[i],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, dataDetail)
	}

	return
}

// GetCommunityPostList 根据community_id 返回该社区id下的所有帖子
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {

	// 1.根据p.CommunityId从redis中获取对应的id列表（key）
	ids, err := redis.GetCommunityPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) returned 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids)) //将结果打印到日志里面，检查id列表获取情况

	// 2.根据id列表，从MySQL中查询帖子的详细信息
	postList, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids)) //将结果打印到日志里面，检查帖子列表的获取情况

	// 查询每篇帖子的 投票数
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 与GetPostList()的处理方式相同，组装结果切片
	//将帖子作者、分区信息 查询出并填充到帖子结果切片中
	data = make([]*models.ApiPostDetail, 0, len(postList)) //为结果数组分配空间

	for i, post := range postList {

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
			VoteNum:         votedata[i],
			Post:            post,
			CommunityDetail: community,
		}

		data = append(data, dataDetail)
	}

	return
}

// GetPostListRouter 根据community参数选择帖子展示的方式
func GetPostListRouter(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 如果ID中带有 communityID, 则根据社区分区进行展示，否则正常展示所有要求范围内的帖子
	if p.CommunityID == 0 {
		data, err = GetPostListInOrder(p)
	} else {
		data, err = GetCommunityPostList(p)
	}

	if err != nil {
		zap.L().Error("GetPostListRouter failed", zap.Error(err))
	}
	return
}
