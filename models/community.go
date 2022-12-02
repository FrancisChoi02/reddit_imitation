package models

import "time"

//社区功能相关的结构体

//使用sqlx库 mysql对应的tag是 db
type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"id" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction" db:"introduction"`
	CreateTime   time.Time `json:"create_Time" db:"create_Time"`
	UpdateTime   time.Time `json:"update_Time" db:"update_time"`
}
