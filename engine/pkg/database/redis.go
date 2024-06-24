package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

// RedisClient 创建一个全局的redis客户端实例Rdb
var RedisClient *redis.Client
var ctx = context.Background()

func InitDefault() {
	InitRedisConn(&RedisOption{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
		PoolSize: 10,
	})
}

func InitRedisConn(option *RedisOption) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: option.Addr,
		//Addr: config.RedisHostPort,
		//留空为没设密码
		Password: option.Addr,
		// 默认的DB 为 0
		DB: option.DB,
		// 连接池大小
		PoolSize: option.PoolSize,
	})
	// 验证是否连接到redis服务端
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("NewRedisConnError:", err)
	}
}

// GiveLike 设置bitmap使用
// 语法: setbit like_id1 100 1   //设置ID为100为1
func GiveLike(keys string, userid int64) (bool, error) {
	res, err := RedisClient.GetBit(ctx, keys, userid-1).Result()
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	res, err = RedisClient.SetBit(ctx, keys, userid-1, 1).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

// 查询bitmap使用
func GiveLikeSelect(keys string, userid int64) (bool, error) {
	res, err := RedisClient.GetBit(ctx, keys, userid-1).Result()
	if err != nil {
		return false, err
	}

	if res == 1 {
		return true, nil
	}

	return false, nil
}

// https://github.com/bsm/redislock
// redis加锁
func Lock(key string) bool {
	var mutex sync.Locker
	mutex.Lock()
	defer mutex.Unlock()
	bool, err := RedisClient.SetNX(ctx, key, 1, 10*time.Second).Result()
	if err != nil {
		log.Println(err.Error())
	}
	return bool
}

// database 释放锁
func UnLock(key string) int64 {
	nums, err := RedisClient.Del(ctx, key).Result()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return nums
}

type RedisOption struct {
	Addr     string
	DB       int
	Password string
	PoolSize int
}
