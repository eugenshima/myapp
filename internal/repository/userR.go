package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eugenshima/myapp/internal/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// UserRedisConnection represents a redis connection
type UserRedisConnection struct {
	rdb *redis.Client
}

// NewUserRedisConnection creates a new connection
func NewUserRedisConnection(rdb *redis.Client) *UserRedisConnection {
	return &UserRedisConnection{rdb: rdb}
}

// Set func inserting entity to redis database
func (rdb *UserRedisConnection) Set(ctx context.Context, user *model.User) error {
	val, err := json.Marshal(model.UserRedis{
		Login:        user.Login,
		Password:     user.Password,
		Role:         user.Role,
		RefreshToken: user.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf(" Marshal: %w", err)
	}
	_, err = rdb.rdb.Set(ctx, user.ID.String(), val, TTL).Result()
	if err != nil {
		return fmt.Errorf(" Set: %w", err)
	}
	return nil
}

// Get func getting entity from redis database
//nolint:dupl
func (rdb *UserRedisConnection) Get(ctx context.Context, id uuid.UUID) (*model.User, error) {
	val, err := rdb.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, fmt.Errorf(" Get: %w", err)
	}
	err = rdb.rdb.Expire(ctx, id.String(), TTL).Err()
	if err != nil {
		return nil, fmt.Errorf(" Expire: %w", err)
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, fmt.Errorf(" Unmarshal: %w", err)
	}
	return user, nil
}

// Delete deletes the user from the cache
func (rdb *UserRedisConnection) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := rdb.rdb.Del(ctx, id.String()).Result()
	if err != nil || res == 0 {
		return fmt.Errorf(" Del: %w", err)
	}
	return nil
}

// GetRefreshToken func gets a refresh token from given user
func (rdb *UserRedisConnection) GetRefreshToken(ctx context.Context, id uuid.UUID) ([]byte, error) {
	val, err := rdb.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		return nil, fmt.Errorf(" Get: %w", err)
	}
	err = rdb.rdb.Expire(ctx, id.String(), TTL).Err()
	if err != nil {
		return nil, fmt.Errorf(" Expire: %w", err)
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, fmt.Errorf(" Unmarshal: %w", err)
	}
	return user.RefreshToken, nil
}

// SetRefreshToken sets the refresh token for the user
func (rdb *UserRedisConnection) SetRefreshToken(ctx context.Context, id uuid.UUID, token []byte) error {
	val, err := rdb.rdb.Get(ctx, id.String()).Result()
	if err != nil {
		return fmt.Errorf(" Get: %w", err)
	}
	user := &model.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return fmt.Errorf(" Unmarshal: %w", err)
	}
	res, err := json.Marshal(model.UserRedis{
		Login:        user.Login,
		Password:     user.Password,
		Role:         user.Role,
		RefreshToken: token,
	})
	if err != nil {
		return fmt.Errorf(" Marshal: %w", err)
	}
	_, err = rdb.rdb.Set(ctx, id.String(), res, TTL).Result()
	if err != nil {
		return fmt.Errorf(" Set: %w", err)
	}
	return nil
}
