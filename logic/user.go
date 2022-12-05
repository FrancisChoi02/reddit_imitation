package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断当前用户是否存在

	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询错误
		return err
	}

	// 2.生成UID
	userID := snowflake.GenID()
	//构建用户实例
	user := models.User{
		Username: p.Username,
		Password: p.Password,
		UserID:   userID,
	}
	//数据保存到MySQL
	return mysql.InsertUser(&user)
}

func Login(p *models.ParamLogin) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是user指针，因此可以获得user.UserID
	if err := mysql.Login(user); err != nil { //从数据库中查找后，user的值会被覆盖一遍
		return "", err
	}

	//生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
