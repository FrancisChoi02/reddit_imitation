package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

// GetCommunityList 查找数据库中的所有community 并返回一个community列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	if err := db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows { //没有数据严格来说不是值得返回的错误
			zap.L().Warn("There is no community in database")
			err = nil
		}
	}
	return
}

//// GetCommunityDetail 根据ID查询社区详情
//func GetCommunityDetail(id int64) (community *models.CommunityDetail, err error) {
//	//新建一个CommunityDetail对象，并为其分配空间
//	community = new(models.CommunityDetail)
//	fmt.Println("The id is", id)
//	sqlStr := `  select community_id,community_name,introduction,create_time
//  				 from community
//  				 where community_id = ?
//				`
//
//	if err := db.Get(&community, sqlStr, id); err != nil {
//		if err == sql.ErrNoRows {
//			fmt.Println("The id is Nothing")
//			err = ErrorInvalidID
//		}
//	}
//	return community, err
//}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select 
			community_id, community_name,introduction,create_Time,update_time
			from community 
			where community_id = ?
	`

	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}

	return community, err
}
