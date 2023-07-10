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

func NewProducer(rdb *redis.Client) *RedisProducer {
	return &RedisProducer{rdb: rdb}
}

func (rdbClient *RedisProducer) RedisProducer() {
	identificator := 0
	for {
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

		time.Sleep(2 * time.Second)
	}
}
