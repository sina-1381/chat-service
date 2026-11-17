package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

// cache logic. this connects to redit port set in .env . inputs key
// and value if value is empty tries to get back the value set based
// on the ket input inside the cache if value is not empty sets the
// ket value inside the cache for 1 hour
func CacheRun(key string, value interface{}) {

	// redis connction logic
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

	if value == "" { // chacks if value is nil to Get the key value
		if err := mycache.Get(ctx, key, &value); err == nil {
			fmt.Println(value)
		}
	} else { // if value not empty it Sets the key value provided inside cache
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
