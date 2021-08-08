/**
 * @Author KYIMH
 * @Description
 * @Date 2021/8/7 15:27
 **/

package redis

import (
	"github.com/go-redis/redis"
	"time"
)

//redis client operators
type Client interface {
	GetClient(redisTag string) (*redis.Client, error)
	Close() error
}

//redis data operators
type Dal interface {
	GetClient(redisTag string) (*redis.Client, error)
	Close() error
	RedisSet(redisTag string, key string, value interface{}, expire int) error
	RedisKeyExists(redisTag string, key string) (bool, error)
	RedisGet(redisTag string, key string) (string, error)
	RedisGetResult(redisTag string, key string) (interface{}, error)
	RedisGetInt(redisTag string, key string) (int, error)
	RedisGetInt64(redisTag string, key string) (int64, error)
	RedisGetUint64(redisTag string, key string) (uint64, error)
	RedisGetFloat64(redisTag string, key string) (float64, error)
	RedisExpire(redisTag string, key string, expire int) error
	RedisPTTL(redisTag string, key string) (int, error)
	RedisTTL(redisTag string, key string) (int, error)
	RedisDel(redisTag string, key string) error
	RedisHGet(redisTag string, key string, field string) (string, error)
	RedisHSet(redisTag string, key string, field string, value string) error
	RedisHDel(redisTag string, key string, field string) error
	RedisZAdd(redisTag string, key string, member string, score string) error
	RedisZRank(redisTag string, key string, member string) (int, error)
	RedisZRange(redisTag string, key string, start int, stop int) (values []string, err error)
	RedisZRangeWithScores(redisTag string, key string, start int, stop int) (values []redis.Z, err error)
	RedisZRem(redisTag string, key string, member string) error
	RedisRPUSH(redisTag string, key string, member string) (err error)
	RedisBLPOP(redisTag string, timeout time.Duration, keys ...string) (value []string, err error)
	RedisLLEN(redisTag string, key string) (value int64, err error)
	RedisLRange(redisTag string, key string, start int, stop int) (values []string, err error)
	RedisKeys(redisTag string, pattern string) (keys []string, err error)
	RedisListAllValuesWithPrefix(redisTag string, prefix string) (map[string]string, error)
	RedisBatchDel(redisTag string, key ...string) error
	RedisMset(redisTag string, pairs ...interface{}) error
	getKeys(redisTag string, prefix string) ([]string, error)
	getKeyAndValuesMap(redisTag string, keys []string, prefix string) (map[string]string, error)
}
