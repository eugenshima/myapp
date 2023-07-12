// Package producer is a producer for redis stream
package producer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisProducer represents the
type RedisProducer struct {
	rdb *redis.Client
}

// NewProducer creates a new RedisProducer
func NewProducer(rdb *redis.Client) *RedisProducer {
	return &RedisProducer{rdb: rdb}
}

// RedisProducer creates a new redis stream entry
func (rdbClient *RedisProducer) RedisProducer(ctx context.Context, stopCh <-chan struct{}) {
	identificator := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Goroutine Closed...")
			return
		case <-stopCh:
			fmt.Println("Goroutine Closed...")
			return
		default:
			id := strconv.FormatInt(time.Now().Unix(), 10)
			payload := map[string]interface{}{
				"timestamp": id,
				"content":   fmt.Sprintf("Redis streaming %d...", identificator),
			}

			identificator++
			id = id + "-" + strconv.Itoa(identificator)

			err := rdbClient.rdb.XAdd(context.Background(), &redis.XAddArgs{
				Stream: "testStream",
				MaxLen: 0,
				ID:     id,
				Values: payload,
			}).Err()
			if err != nil {
				fmt.Println("Error adding message to Redis Stream:", err)
			}
			const TTL = 2
			time.Sleep(TTL * time.Second)
		}
	}

}
