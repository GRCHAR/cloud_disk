package redis

import "cloud_disk/src/logger"

type distributedLock struct {
	isLock bool
}

func init() {

}

func CreateDistributedLock(key string) bool {
	_, err := Rdb.Set(key, false, 1000).Result()
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	return true

}

func DeleteDistributedLock(key string) bool {
	_, err := Rdb.Del(key).Result()
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	return true
}

func GetDistributedLock(key string) bool {
	_, err := Rdb.Set(key, true, 1000).Result()
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	return true
}

func FreeDistributedLock(key string) bool {
	_, err := Rdb.Set(key, false, 1000).Result()
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	return true
}

func FindDistributedLock(key string) (bool, error) {
	result, err := Rdb.Get(key).Result()
	if err != nil {
		logger.GetLogger().Error(err)
		return false, err
	}
	if result == "true" {
		return true, nil
	}
	return false, nil
}
