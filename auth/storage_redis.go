package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type storageRedis struct {
	redisClient *redis.Client
}

func (s *storageRedis) Get(prefix, key string) ([]byte, error) {
	data, err := s.redisClient.Get(context.Background(), fmt.Sprintf("session_%s_%s", prefix, key)).Bytes()
	if err == redis.Nil {
		return nil, ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *storageRedis) Set(prefix, key string, data []byte, duration time.Duration) error {
	err := s.redisClient.Set(context.Background(), fmt.Sprintf("session_%s_%s", prefix, key), string(data), duration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *storageRedis) Delete(prefix, key string) error {
	err := s.redisClient.Del(context.Background(), fmt.Sprintf("session_%s_%s", prefix, key)).Err()
	if err != nil {
		return err
	}

	return nil
}

func initStorageRedis(args ...string) (storage, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("redis storage require at least connection string as argument for initialization")
	}

	var password string
	if len(args) >= 2 {
		password = args[1]
	}

	newStorage := &storageRedis{}

	newStorage.redisClient = redis.NewClient(&redis.Options{
		Addr:     args[0],
		Password: password,
		DB:       0, // use default DB
	})

	// Pinging here can cause panic if docker decides to start our service faster than Redis.
	/*_, err := newStorage.redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}*/

	return newStorage, nil
}

func init() {
	storages = append(storages, storageInitialization{
		StorageType: "redis",
		Initialize:  initStorageRedis,
	})
}
