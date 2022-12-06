package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// getIDSFormKey 根据size和page参数 从redis中从大到小获取帖子id
func getIDSFormKey(key string, page, size int64) ([]string, error) {
	//确定索引的起始点
	start := (page - 1) * size
	end := start + (size - 1)

	// Zrevrange 按分数从大到小的顺序 查询
	return client.ZRevRange(key, start, end).Result()
}

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

func GetPostVoteData(ids []string) (data []int64, err error) {

	//将所有需要搜寻id的数据 集合一起发给redis，节省每个命令的网络往返时间（RTT）
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id) // id帖子对应的  用户 -- 投票 键值对
		pipeline.ZCount(key, "1", "1")                  //只统计赞成票 "1" 的数量
	}
	cmders, err := pipeline.Exec() //cmders 是帖子投票情况切片
	if err != nil {
		return nil, err
	}

	//定义结果数组并为其分配空间
	data = make([]int64, 0, len(ids))

	//cmder是每一个每一个id帖子 用户投票的总数，类型是Cmder 需要转换
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val() //转换
		data = append(data, v)
	}

	return
}

func GetCommunityPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 获取排序形式order参数
	order := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		order = getRedisKey(KeyPostScoreZSet)
	}

	//使用 zinterstore 将社区分区中的帖子set 与 帖子分数的zset 生成一个心得zset
	//用之前的逻辑 从新的zset中获取数据

	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存key减少 zinterstore 执行的次数
	// E.g newKey = "post:score1" ; order="post:score"; communityID="1"
	// newkey是用于保存 order排序下，社区为id的 临时Zset
	newKey := order + strconv.Itoa(int(p.CommunityID))
	if client.Exists(newKey).Val() < 1 {
		// 不存在这个社区对应的键值对，需要计算
		pipeline := client.Pipeline()

		//newKey是 cKey 和 order 的交集，结果集的score为最大值
		pipeline.ZInterStore(newKey, redis.ZStore{
			Aggregate: "MAX",
		}, cKey, order)

		//缓存key，设置过期时间，60秒后有访问再更新一次，减少 zinterstore 执行的次数
		pipeline.Expire(newKey, 60*time.Second) //设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	//newKey是一个新的Zset，带过期时间
	return getIDSFormKey(newKey, p.Page, p.Size)

}
