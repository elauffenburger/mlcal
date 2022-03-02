package calcache

import (
	"encoding/json"

	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/go-redis/redis"
)

const redisCalKey = "cal_ics"

type redisCalValue struct {
	Calendar mlcal.Calendar
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client}
}

func (s *RedisCache) Get() (*mlcal.Calendar, error) {
	calValueStr, err := s.client.Get(redisCalKey).Result()
	if err != nil {
		return nil, err
	}

	var cal redisCalValue
	err = json.Unmarshal([]byte(calValueStr), &cal)
	if err != nil {
		return nil, err
	}

	return &cal.Calendar, nil
}

func (s *RedisCache) GetICS() (string, error) {
	cal, err := s.Get()
	if err != nil {
		return "", err
	}

	return cal.ToICS().Serialize(), nil
}

func (s *RedisCache) Set(cal *mlcal.Calendar) error {
	calValue := redisCalValue{*cal}
	calValueBytes, err := json.Marshal(&calValue)
	if err != nil {
		return err
	}

	_, err = s.client.Set(redisCalKey, string(calValueBytes), 0).Result()

	return err
}
