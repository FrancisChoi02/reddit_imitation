package models

//请求参数结构体

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
