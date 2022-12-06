package models

//请求参数结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParmVoteData 投票参数
type ParmVoteData struct {
	// UserID 可以从当前的上下文的用户信息中获取
	PostID    string `json:"post_id" binding:"required"`             //前端有可能传来字符串，为反序列化作准备
	Direction int    `json:"direction,string" biding:"oneof=1 0 -1"` //赞成（1） 或者 反对票（-1；
	// 不需要required 不然会把0值过滤了
}

// ParamPostList 获取帖子列表URL中的query string参数
type ParamPostList struct {
	Page        int64  `json:"page" form:"page"` //因为从URL中获取参数，因此使用form tag
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"` //可以为空 初值是0
}
