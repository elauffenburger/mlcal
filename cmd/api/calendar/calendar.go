package calendar

import (
	"fmt"
	"log"
	"os"

	"github.com/elauffenburger/musical-literacy-cal/cmd/api/calendar/calcache"
	"github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"
	"github.com/go-redis/redis"
)

type Getter interface {
	Get() (*mlcal.Calendar, error)
	GetICS() (string, error)
}

type Setter interface {
	Set(*mlcal.Calendar) error
}

type Cache interface {
	Getter
	Setter
}

func CreateRedisCache(redisOpts *redis.Options) (*calcache.RedisCache, error) {
	logger := log.New(os.Stdout, "[redis-cache] ", log.LstdFlags)
	logger.Printf("using redis addr: '%s'", redisOpts.Addr)

	// Create a redis client and do a quick connection check.
	client := redis.NewClient(redisOpts)
	if err := client.Get("").Err(); err != redis.Nil {
		err = fmt.Errorf("failed to create redis client: %w", err)
		logger.Print(err)

		return nil, err
	}

	return calcache.NewRedisCache(client), nil
}
