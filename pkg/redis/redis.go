package redis

import (
	"context"
	"os"
	"sync"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *goredis.Client
)

func Init() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil
	}

	var initErr error
	once.Do(func() {
		client = goredis.NewClient(&goredis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
		})
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := client.Ping(ctx).Err(); err != nil {
			initErr = err
			client = nil
		}
	})
	return initErr
}

func Client() *goredis.Client {
	return client
}

func Enabled() bool {
	return client != nil
}
