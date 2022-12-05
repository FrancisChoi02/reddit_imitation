package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

//将数据库的相关操作都封装成函数
//待logic逻辑层进行调用

const key = "Francis"

// CheckUserExist 查找数据库中用户是否存在
func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username =?`

	//查找数据库中，有多少个该用户名的数据
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}

	if count > 0 {
		//用户已经存在，不需要重复注册
		return ErrorUserExist
	}
	return nil
}

func encryptPassword(oPassword string) string {
	//用key进行加密
	h := md5.New()
	h.Write([]byte(key))
	tmp := h.Sum([]byte(oPassword))

	return hex.EncodeToString(tmp)
}

// InsertUser 将新用户记录到数据库中
func InsertUser(user *models.User) (err error) {
	//执行SQL语句入库
	sqlStr := `insert into user(user_id,username,password) value(?,?,?)`

	//数据库不可以保存密码明文
	user.Password = encryptPassword(user.Password)
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username= ?`

	err = db.Get(user, sqlStr, user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		return err
	}

	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(pid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id = ?`
	err = db.Get(user, sqlStr, pid)
	return
}
