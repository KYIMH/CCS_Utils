package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"
)

type ClientPoolType map[string]*redis.Client

type ClientImpl struct {
	Pool ClientPoolType
}

//create new redis client
func NewClient() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		//redis config
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,

		//connection pool
		PoolSize:     15, // socket connect nums
		MinIdleConns: 10, // Idle connect nums

		//redis client io timeouts
		DialTimeout:  5 * time.Second, //max time to connect redis
		ReadTimeout:  3 * time.Second, //max time of read
		WriteTimeout: 3 * time.Second, //max time of write
		PoolTimeout:  4 * time.Second, //max wait time

		//idle connection check, include IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //frequency of idle check
		IdleTimeout:        5 * time.Minute,  //max time for a idle connection
		MaxConnAge:         0 * time.Second,  //connection live time

		//strategies when command failed
		MaxRetries:      0,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,

		Dialer: func() (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial("tcp", "127.0.0.1:6379")
		},

		//hook
		OnConnect: func(conn *redis.Conn) error {
			fmt.Printf("conn=%v\n", conn)
			return nil
		},
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		logrus.Error("redis connection failed: ", err.Error())
	}

	return redisClient
}

//get redis client by redis tag name
func (c ClientImpl) GetClient(redisTag string) (*redis.Client, error) {

	cli, ok := c.Pool[redisTag]

	if !ok {
		return nil, errors.New("no connection " + redisTag + " in Manager")
	}

	return cli, nil
}

//close all redis client
func (c ClientImpl) Close() error {

	for _, cli := range c.Pool {
		err := cli.Close()
		if nil != err {
			return err
		}
	}

	return nil
}

//redis String set
func (c ClientImpl) RedisSet(redisTag string, key string, value interface{}, expire int) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	if expire > 0 {
		err := cli.Do("SET", key, value, "PX", expire).Err()
		if err != nil {
			logrus.Error("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := cli.Do("SET", key, value).Err()
		if err != nil {
			logrus.Error("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func (c ClientImpl) RedisKeyExists(redisTag string, key string) (bool, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return false, err
	}

	ok, err := cli.Do("EXISTS", key).Bool()

	return ok, err
}

func (c ClientImpl) RedisGet(redisTag string, key string) (string, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return "", err
	}

	value, err := cli.Do("GET", key).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func (c ClientImpl) RedisGetResult(redisTag string, key string) (interface{}, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return nil, err
	}

	v, err := cli.Do("GET", key).Result()
	if err == redis.Nil {
		return v, nil
	}

	return v, err
}

func (c ClientImpl) RedisGetInt(redisTag string, key string) (int, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return 0, err
	}

	v, err := cli.Do("GET", key).Int()
	if err == redis.Nil {
		return 0, nil
	}

	return v, err
}

func (c ClientImpl) RedisGetInt64(redisTag string, key string) (int64, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return 0, err
	}

	v, err := cli.Do("GET", key).Int64()
	if err == redis.Nil {
		return 0, nil
	}

	return v, err
}

func (c ClientImpl) RedisGetUint64(redisTag string, key string) (uint64, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return 0, err
	}

	v, err := cli.Do("GET", key).Uint64()
	if err == redis.Nil {
		return 0, nil
	}

	return v, err
}

func (c ClientImpl) RedisGetFloat64(redisTag string, key string) (float64, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return 0.0, err
	}

	v, err := cli.Do("GET", key).Float64()
	if err == redis.Nil {
		return 0.0, nil
	}

	return v, err
}

func (c ClientImpl) RedisExpire(redisTag string, key string, expire int) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("EXPIRE", key, expire).Err()
	if err != nil {
		logrus.Error("RedisExpire Error!", key, "Details:", err.Error())
		return err
	}

	return nil
}

func (c ClientImpl) RedisPTTL(redisTag string, key string) (int, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return -1, err
	}

	ttl, err := cli.Do("PTTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func (c ClientImpl) RedisTTL(redisTag string, key string) (int, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return -1, err
	}

	ttl, err := cli.Do("TTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func (c ClientImpl) RedisDel(redisTag string, key string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("DEL", key).Err()
	if err != nil {
		logrus.Error("RedisDel Error! key:", key, "Details:", err.Error())
	}

	return err
}

func (c ClientImpl) RedisHGet(redisTag, key, field string) (string, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return "", err
	}

	value, err := cli.Do("HGET", key, field).String()
	if err != nil {
		logrus.Error("HGet Error! key:", key)
	}

	return value, nil
}

func (c ClientImpl) RedisHSet(redisTag, key, field, value string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("HSET", key, field, value).Err()
	if err != nil {
		logrus.Error("RedisHSet Error!", key, "field:", field, "Details:", err.Error())
	}

	return err
}

func (c ClientImpl) RedisHDel(redisTag, key, field string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("HDEL", key, field).Err()
	if err != nil {
		logrus.Error("RedisHDel Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func (c ClientImpl) RedisZAdd(redisTag, key, member, score string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("ZADD", key, score, member).Err()
	if err != nil {
		logrus.Error("RedisZAdd Error!", key, "member:", member, "score:", score, "Details:", err.Error())
	}
	return err
}

func (c ClientImpl) RedisZRank(redisTag, key, member string) (int, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return -1, err
	}

	rank, err := cli.Do("ZRANK", key, member).Int()
	if err == redis.Nil {
		return -1, nil
	}

	if err != nil {
		logrus.Error("RedisZRank Error!", key, "member:", member, "Details:", err.Error())
		return -1, nil
	}

	return rank, err
}

func (c ClientImpl) RedisZRange(redisTag string, key string, start, stop int) (values []string, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []string{}, err
	}

	values, err = cli.ZRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		logrus.Error("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisZRangeWithScores(redisTag string, key string, start, stop int) (values []redis.Z, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []redis.Z{}, err
	}

	values, err = cli.ZRangeWithScores(key, int64(start), int64(stop)).Result()
	if err != nil {
		logrus.Error("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisZRem(redisTag, key, member string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("ZREM", key, member).Err()
	if err != nil {
		logrus.Error("RedisZRem Error!", key, "member:", member, "Details:", err.Error())
	}
	return err
}

func (c ClientImpl) RedisRPUSH(redisTag string, key string, member string) (err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Do("RPUSH", key, member).Err()
	if err != nil {
		logrus.Error("RedisRPUSH Error!", key, member, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisBLPOP(redisTag string, timeout time.Duration, keys ...string) (value []string, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []string{}, err
	}

	value, err = cli.BLPop(timeout, keys...).Result()
	if err == redis.Nil {
		err = nil
		return
	}

	if err != nil {
		logrus.Error("BLPop Error!", keys, timeout, "Details:", err.Error())
		return
	}
	return
}

func (c ClientImpl) RedisLLEN(redisTag string, key string) (value int64, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return 0, err
	}

	value, err = cli.LLen(key).Result()
	if err != nil {
		logrus.Error("RedisLLEN Error!", key, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisLRange(redisTag string, key string, start, stop int) (values []string, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []string{}, err
	}

	values, err = cli.LRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		logrus.Error("RedisLRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisKeys(redisTag string, pattern string) (keys []string, err error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []string{}, err
	}

	keys, err = cli.Keys(pattern).Result()
	if err != nil {
		logrus.Error("RedisKeys Error!", pattern, "Details:", err.Error())
		return
	}

	return
}

func (c ClientImpl) RedisListAllValuesWithPrefix(redisTag string, prefix string) (map[string]string, error) {

	keys, err := c.getKeys(redisTag, fmt.Sprintf("%s*", prefix))
	if err != nil {
		return nil, err
	}

	values, err := c.getKeyAndValuesMap(redisTag, keys, prefix)

	return values, nil
}

func (c ClientImpl) RedisBatchDel(redisTag string, key ...string) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.Del(key...).Err()
	if err != nil {
		logrus.Error("RedisBatchDel Error! key:", key, "Details:", err.Error())
	}

	return err
}

func (c ClientImpl) RedisMset(redisTag string, pairs ...interface{}) error {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return err
	}

	err = cli.MSet(pairs...).Err()
	if err != nil {
		logrus.Error("RedisMset Error! pairs:", pairs, "Details:", err.Error())
	}

	return err
}

func (c ClientImpl) getKeys(redisTag string, prefix string) ([]string, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return []string{}, err
	}

	var allKeys []string
	var cursor uint64
	count := int64(10)

	for {
		var keys []string
		var err error
		keys, cursor, err = cli.Scan(cursor, prefix, count).Result()
		if err != nil {
			return nil, nil
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}

	}

	return allKeys, nil
}

func (c ClientImpl) getKeyAndValuesMap(redisTag string, keys []string, prefix string) (map[string]string, error) {

	cli, err := c.GetClient(redisTag)
	if nil != err {
		return nil, err
	}

	values := make(map[string]string)
	for _, key := range keys {
		value, err := cli.Do("GET", key).String()
		if err != nil {
			logrus.Error("error retrieving value for key ", key, "Details:", err.Error())

		}

		strippedKey := strings.Split(key, prefix)
		values[strippedKey[1]] = value
	}

	return values, nil
}
