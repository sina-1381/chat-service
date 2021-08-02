package services

import (
	"context"
	"fmt"
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"os"
	"time"
)

func CacheRun(key string, value interface{}) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": os.Getenv("CACHE_PORT"),
		},
	})
	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	ctx := context.TODO()
	if value == "" {
		if err := mycache.Get(ctx, key, &value); err == nil {
			fmt.Println(value)
		}
	} else {
		if err := mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   key,
			Value: value,
			TTL:   1 * time.Hour,
		}); err != nil {
			panic(err)
		}
	}
}
