package mysql

import "bluebell/models"

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
   	post_id,title,content, author_id, community_id)              
	value (?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
    			post_id, title, content, author_id, community_id, create_time
				from post 
				where post_id = ?`

	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostList(page int64, size int64) (postList []*models.Post, err error) {
	sqlStr := `select 
    			post_id, title, content, author_id, community_id, create_time
				from post  
				limit ?,?
				` //要有条数限制  为后续分页做准备

	postList = make([]*models.Post, 0, size)
	//需要传地址才能被修改覆盖
	err = db.Select(&postList, sqlStr, (page-1)*size, size) //从(page - 1)*size个后开始取，取size个
	return postList, err
}
