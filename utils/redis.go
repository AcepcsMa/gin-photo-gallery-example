package utils

import (
	"fmt"
	"gin-photo-storage/conf"
	"gin-photo-storage/constant"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var RedisClient *redis.Client
var InitComplete = make(chan struct{}, 1)

// Init redis client
func init() {
	host := conf.ServerCfg.Get(constant.REDIS_HOST)
	port := conf.ServerCfg.Get(constant.REDIS_PORT)
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB: 0,
	})
	InitComplete <- struct{}{}
}

// Add an auth to redis, meaning that he/she has logged in.
func AddAuthToRedis(username string) error {
	key := fmt.Sprintf("%s%s", constant.LOGIN_USER, username)
	err := RedisClient.Set(key, username, constant.LOGIN_MAX_AGE * time.Second).Err()
	if err != nil {
		//log.Fatalln(err)
		AppLogger.Info(err.Error(), zap.String("service", "AddAuthToRedis()"))
		return err
	}
	return nil
}

// Check if an auth is in redis.
func IsAuthInRedis(username string) bool {
	key := fmt.Sprintf("%s%s", constant.LOGIN_USER, username)
	err := RedisClient.Get(key).Err()
	if err != nil {
		//log.Println(err)
		AppLogger.Info(err.Error(), zap.String("service", "IsAuthInRedis()"))
		return false
	}
	return true
}

// Remove an auth from redis, meaning he/she is logging out.
func RemoveAuthFromRedis(username string) bool {
	key := fmt.Sprintf("%s%s", constant.LOGIN_USER, username)
	err := RedisClient.Del(key).Err()
	if err != nil {
		//log.Println(err)
		AppLogger.Info(err.Error(), zap.String("service", "RemoveAuthFromRedis()"))
		return false
	}
	return true
}

// Set the upload status for a photo.
func SetUploadStatus(key string, value int) bool {
	err := RedisClient.Set(key, value, 0).Err()
	if err != nil {
		//log.Println(err)
		AppLogger.Info(err.Error(), zap.String("service", "SetUploadStatus()"))
		return false
	}
	return true
}

// Get the upload status of a photo.
func GetUploadStatus(key string) int {
	val := RedisClient.Get(key).Val()
	if val == "" {
		return -2	// means no such key
	}
	status, _ := strconv.Atoi(val)
	return status
}

// Send a message to the given channel.
func SendToChannel(channel string, message string) bool {
	err := RedisClient.Publish(channel, message).Err()
	if err != nil {
		//log.Println(err)
		AppLogger.Info(err.Error(), zap.String("service", "SendToChannel()"))
		return false
	}
	return true
}