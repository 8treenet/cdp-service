package repository

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

var defExpiration time.Duration = time.Minute * 10

func redisJSONGet(client redis.Cmdable, key string, value interface{}) error {
	bytes, err := client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

func redisJSONSet(client redis.Cmdable, key string, value interface{}, expirations ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiration := defExpiration
	if len(expirations) > 0 {
		expiration = expirations[0]
	}
	return client.Set(key, data, expiration).Err()
}
