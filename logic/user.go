package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
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

func Login(p *models.ParamSignUp) (err error) {
	user := models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(&user)
}
