package strategies

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	_ LimiterStrategyInterface = &RedisLimiterStrategy{}
)

const (
	KeyWithoutTTL = -1
	KeyNotFound   = -2
)

type RedisLimiterStrategy struct {
	Client *redis.Client
	Logger zerolog.Logger
	Now    func() time.Time
}

func NewRedisLimiterStrategy(
	client *redis.Client,
	logger zerolog.Logger,
	now func() time.Time,
) *RedisLimiterStrategy {
	return &RedisLimiterStrategy{
		Client: client,
		Logger: logger,
		Now:    now,
	}
}

func (rls *RedisLimiterStrategy) Check(ctx context.Context, r *RateLimiterRequest) (*RateLimiterResult, error) {
	p := rls.Client.Pipeline()
	getResult := p.Get(ctx, r.Key)
	ttlResult := p.TTL(ctx, r.Key)

	if _, err := p.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var ttlDuration time.Duration

	ttl, err := ttlResult.Result()
	if err != nil || ttl == KeyWithoutTTL || ttl == KeyNotFound {
		ttlDuration = r.Duration

		if err := rls.Client.Expire(ctx, r.Key, r.Duration).Err(); err != nil {
			return nil, err
		}
	} else {
		ttlDuration = ttl
	}

	currentCount, err := getResult.Int64()
	if err != nil && errors.Is(err, redis.Nil) {
		// Fail-safe in case there's an error while getting the count
		currentCount = 0
	}

	expiresAt := rls.Now().Add(ttlDuration)

	if currentCount >= r.Limit {
		return &RateLimiterResult{
			Result:    Deny,
			Total:     currentCount,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	incrResult := rls.Client.Incr(ctx, r.Key)
	nextTotal, err := incrResult.Result()
	if err != nil {
		return nil, err
	}

	if nextTotal > r.Limit {
		return &RateLimiterResult{
			Result:    Deny,
			Total:     nextTotal,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	return &RateLimiterResult{
		Result:    Allow,
		Total:     nextTotal,
		Limit:     r.Limit,
		Remaining: r.Limit - nextTotal,
		ExpiresAt: expiresAt,
	}, nil
}
