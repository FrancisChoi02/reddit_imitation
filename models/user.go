package models

type User struct {
	//tag对应数据表中的名称
	UserID   int64  `db:"user_id"` //雪花算法生成的ID是int64
	Username string `db:"username"`
	Password string `db:"password"`
}
