package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"time"
)

// 帖子投票分数算法：
/*
投票情况分类:
direction = 1:
之前没投票 		+432   （1）
之前投反对票		+432 *2   （2）

direction = 0
之前投反对票		+432   （1）
之前投赞成票		-432   （1）

direction = -1
之前没投票		-432   （1）
之前投赞成票		-432 *2   （2）

投票限制：
1.发帖后一个星期不允许用户对该帖子投票（避免挖坟）
2.到期后将redis 保存赞成票数和反对票数的键值对 KeyPostScoreZSet 保存到MySQL（持久化）
3.到期后 删除 redis中KeyPostVotedZSetPrefix的赞踩选项，用户无法再投票
*/

// 投一票就加 432分 投满200赞成票就维持在榜上一天 --> 86400/200 = 432  24h * 60min * 60s = 86400
const (
	secondsInOneWeek = 7 * 24 * 60 * 60
	scorePerVote     = 432 //每一票代表的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID int64) (err error) {
	//两个键值对的保存，要保证全部都能正常执行，因此要用 事务
	pipeline := client.TxPipeline()

	//将当前帖子的创建时间 记录到redis的  帖子--时间 键值对中
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()), //当前时间（秒）
		Member: postID,
	})

	//将当前帖子的初始分数， 记录到redis的 帖子--分数  键值对中
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()), //当前时间（秒）
		Member: postID,
	})

	_, err = pipeline.Exec()
	return err
}

// VoteForPost 为帖子进行投票
func VoteForPost(userID, postID string, direction float64) (err error) { //direction转变为float64类型，因为
	// 1.判断投票限制
	// 查看帖子发布时间（mysql.CreatePost的时候, Redis要记录）
	postime := client.ZScore(getRedisKey(KeyPostScoreZSet), postID).Val()
	if float64(time.Now().Unix())-postime > secondsInOneWeek { //当前时间（秒）与帖子发布时间的差值
		return ErrVoteTimeExpire
	}

	// 2.更新帖子分数
	//查看当前用户 给 当前帖子 的投票记录
	recordDir := client.ZScore(getRedisKey(KeyPostVotedZSetPrefix+postID), userID).Val()
	if direction == recordDir { //不允许重复投票
		return ErrVoteRepeated
	}

	var dir float64
	if direction > recordDir {
		dir = 1
	} else {
		dir = -1
	}
	dirDif := math.Abs(recordDir - direction) //计算投票记录 和 当前投票决定 横跨选项的差值

	//因为 帖子--分数 键值对和 用户--投票 键值对 要同时执行，所以要使用事务
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), dirDif*dir*scorePerVote, postID)

	// 3.记录 用户 为该帖子投票的数据
	if direction == 0 { //撤销投票
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPrefix+postID), userID)
	} else {
		//Zset会自动去重，会覆盖原有的记录
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPrefix+postID), redis.Z{
			Score:  direction, // 赞成票还是反对票
			Member: userID,
		})
	}
	_, err = pipeline.Exec()
	return err
}
