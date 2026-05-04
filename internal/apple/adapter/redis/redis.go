package redis

import (
	"github.com/redis/go-redis/v9"
	"time"
)

const ttl = time.Hour

type Redis struct {
	client *redis.Client
}

func New(client *redis.Client) *Redis {
	return &Redis{client: client}
}
