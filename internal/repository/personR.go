package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eugenshima/myapp/internal/model"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	rdb *redis.Client
}

func NewRedisConnection(rdb *redis.Client) *RedisConnection {
	return &RedisConnection{rdb: rdb}
}

func (rdb *RedisConnection) RedisGetByID(ctx context.Context, id uuid.UUID) *model.Person {
	val, err := rdb.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		fmt.Printf("failed to get: %v", err)
		return nil
	}
	err = rdb.rdb.Expire(context.Background(), id.String(), 20*time.Minute).Err()
	if err != nil {
		fmt.Println(err)
	}
	person := model.Person{}

	person.ID = id
	if err != nil {
		fmt.Printf("failed to parse: %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(val), &person)
	if err != nil {
		fmt.Printf("failed to unmarshal: %v", err)
		return nil
	}
	return &person
}

func (rdb *RedisConnection) XRedisSetByID(ctx context.Context, entity *model.Person) error {
	streamID, err := rdb.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "Stream1",
		Values: map[string]interface{}{
			"name":      entity.Name,
			"age":       entity.Age,
			"ishealthy": entity.IsHealthy,
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to set: %v", err)
	}
	fmt.Println(streamID)
	return nil
}

func (rdb *RedisConnection) RedisSetByID(ctx context.Context, entity *model.Person) error {
	val, err := json.Marshal(model.PersonRedis{Name: entity.Name, Age: entity.Age, IsHealthy: entity.IsHealthy})
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}
	_, err = rdb.rdb.Set(ctx, entity.ID.String(), val, 20*time.Minute).Result()
	if err != nil {
		return fmt.Errorf("failed to set: %v", err)
	}
	return nil
}
