package nredis

import (
	"time"

	"github.com/go-redis/redis"
)

type Store struct {
	client *redis.Client
}

func NewStore(addr string, password string, index int) *Store {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       index,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &Store{client: client}
}

func (s *Store) Set(key string, val []byte, expire int64) error {
	return s.client.Set(key, val, time.Duration(expire)).Err()
}

func (s *Store) Get(key string) (string, error) {
	return s.client.Get(key).Result()
}

func (s *Store) Expire(key string, expire int64) error {
	return s.client.Expire(key, time.Duration(expire)).Err()
}
