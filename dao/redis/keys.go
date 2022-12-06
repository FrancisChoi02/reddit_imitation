package redis

// redis key 要适当使用前缀字段，方便拆分和查询
// 字段的拆分 用 ： 或者 |

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZSet        = "post:time"   // zset; 帖子 以及 发帖时间
	KeyPostScoreZSet       = "post:score"  // zset; 帖子 以及	帖子累计的投票分数
	KeyPostVotedZSetPrefix = "post:voted:" // zset;	前缀，后面还要加上帖子ID； 记录 用户 以及 投票类型（赞/踩）；//避免重复投票
	KeyCommunitySetPrefix  = "community:"  // set ; 保存每个分区下的帖子id
)

// getRedisKey 为 Redis key 添加前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
