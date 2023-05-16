package redis

import (
	"cloud_disk/src/logger"
	"github.com/go-redis/redis"
	"log"
)

var Rdb = NewRedisManager()

func NewRedisManager() *redis.Client {
	return Rdb
}

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Rdb.Ping().Result()
	if err != nil {
		panic("redis连接失败:" + err.Error())
	}
	log.Println("Connected to redis!")
}

func FindSessionBySessionKey(keyName string) (string, bool) {
	result, err := Rdb.HGet("session", keyName).Result()
	if err != nil {
		logger.GetLogger().Error("find key error:%+v", err)
		return "", false
	}
	return result, true
}

func CreateSession(keyName string, value interface{}) (bool, error) {
	result, err := Rdb.HSet("session", keyName, value).Result()
	if err != nil {
		logger.GetLogger().Error("Create redis failed:", err)
		return false, err
	}
	return result, err
}

func ExistSession(keyName string) bool {
	result, err := Rdb.HExists("session", keyName).Result()
	if err != nil {
		logger.GetLogger().Error("Exist redis failed:", err)
		return false
	}
	return result
}

func GetSessionValue(sessionKey string, keyName string) (string, bool) {
	result, err := Rdb.HGet(sessionKey, keyName).Result()
	if err != nil {
		logger.GetLogger().Error("find key error:%+v", err)
		return "", false
	}
	return result, true
}

func CreateSessionValue(sessionKey string, keyName string, value string) bool {
	_, err := Rdb.HSet(sessionKey, keyName, value).Result()
	if err != nil {
		logger.GetLogger().Error("create key error:%+v", err)
		return false
	}
	return true
}

func FindValueByKey(key string) (string, bool) {
	result, err := Rdb.Get(key).Result()
	if err != nil {
		logger.GetLogger().Error("find key error:%+v", err)
		return "", false
	}
	return result, true
}

func DeleteSessionValue(sessionKey string, keyName string) bool {
	_, err := Rdb.HDel(sessionKey, keyName).Result()
	if err != nil {
		logger.GetLogger().Error("Delete redis failed:", err)
		return false
	}
	return true
}

func DeleteSession(sessionKey string) bool {
	_, err := Rdb.Del(sessionKey).Result()
	if err != nil {
		logger.GetLogger().Error("Delete redis failed:", err)
		return false
	}
	return true
}
