// Package consumer is for redis stream consumer
package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConsumer is a struct for Redis Stream Consumer
type RedisConsumer struct {
	rdb *redis.Client
}

// NewConsumer creates a new Redis Stream Consumer
func NewConsumer(rdb *redis.Client) *RedisConsumer {
	return &RedisConsumer{rdb: rdb}
}

// RedisConsumer reading streams from redis
func (rdbClient *RedisConsumer) RedisConsumer() {
	for {
		streams, err := rdbClient.rdb.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"testStream", "0"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			fmt.Println("Error reading messages from Redis Stream:", err)
		}
		for _, stream := range streams {
			streamName := stream.Stream
			messages := stream.Messages

			for _, msg := range messages {
				messageID := msg.ID
				messageData := msg.Values

				fmt.Println("Received message from Redis Stream:", messageID, messageData)

				_, err := rdbClient.rdb.XDel(context.Background(), streamName, messageID).Result()
				if err != nil {
					fmt.Println("Error deleting message from Redis Stream:", err)
				}
			}
		}
		const TTL = 2
		time.Sleep(TTL * time.Second)
	}
}
