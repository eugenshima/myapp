// Package repository provides functions for interacting with a database
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

// RedisConnection represents a redis connection
type RedisConnection struct {
	rdb *redis.Client
}

// NewRedisConnection creates a new connection
func NewRedisConnection(rdb *redis.Client) *RedisConnection {
	return &RedisConnection{rdb: rdb}
}

// const for Redis
const (
	TTL = 20 * time.Minute
)

// RedisGetByID func returns a Redis entity by ID
//nolint:dupl
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
	err = json.Unmarshal([]byte(val), &person)
	if err != nil {
		return nil, fmt.Errorf(" Unmarshal: %w", err)
	}
	return person, nil
}

// RedisSetByID func inserting entity to redis database
func (rdb *RedisConnection) RedisSetByID(ctx context.Context, entity *model.Person) error {
	val, err := json.Marshal(model.PersonRedis{
		Name:      entity.Name,
		Age:       int(entity.Age),
		IsHealthy: entity.IsHealthy,
	})
	if err != nil {
		return fmt.Errorf(" Marshal: %w", err)
	}
	_, err = rdb.rdb.Set(ctx, entity.ID.String(), val, TTL).Result()
	if err != nil {
		return fmt.Errorf(" Set: %w", err)
	}
	return nil
}

// RedisDeleteByID func deleting entity from redis database
func (rdb *RedisConnection) RedisDeleteByID(ctx context.Context, id uuid.UUID) error {
	res, err := rdb.rdb.Del(ctx, id.String()).Result()
	if err != nil || res == 0 {
		return fmt.Errorf(" Del: %w", err)
	}
	return nil
}
