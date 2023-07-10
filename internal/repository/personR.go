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

const (
	TTL = 20 * time.Minute
)

func (rdb *RedisConnection) RedisGetByID(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	val, err := rdb.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, fmt.Errorf(" Get: %w", err)
	}
	err = rdb.rdb.Expire(ctx, id.String(), TTL).Err()
	if err != nil {
		return nil, fmt.Errorf(" Expire: %w", err)
	}
	person := &model.Person{}

	person.ID = id
	err = json.Unmarshal([]byte(val), &person)
	if err != nil {
		return nil, fmt.Errorf(" Unmarshal: %w", err)
	}
	return person, nil
}

func (rdb *RedisConnection) RedisSetByID(ctx context.Context, entity *model.Person) error {
	val, err := json.Marshal(model.PersonRedis{Name: entity.Name, Age: entity.Age, IsHealthy: entity.IsHealthy})
	if err != nil {
		return fmt.Errorf(" Marshal: %w", err)
	}
	_, err = rdb.rdb.Set(ctx, entity.ID.String(), val, TTL).Result()
	if err != nil {
		return fmt.Errorf(" Set: %w", err)
	}
	return nil
}

func (rdb *RedisConnection) RedisDeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := rdb.rdb.Del(ctx, id.String()).Result()
	if err != nil {
		return fmt.Errorf(" Del: %w", err)
	}
	return nil
}
