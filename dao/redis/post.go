package redis

import "bluebell/models"

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从 Redis中获取ID
	//1. 根据用户请求中携带的 order参数 确定要查询的redis key
	var key string
	if p.Order == models.OrderTime {
		key = getRedisKey(KeyPostTimeZSet)
	} else {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 2. 确定索引的起始点
	start := (p.Page - 1) * p.Size
	end := start + (p.Size - 1)

	// 3. Zrevrange 按分数从大到小的顺序 查询
	return client.ZRevRange(key, start, end).Result()
}
