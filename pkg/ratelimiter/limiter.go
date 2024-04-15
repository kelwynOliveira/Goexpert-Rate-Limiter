package ratelimiter

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	rip "github.com/vikram1565/request-ip"

	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/logger"
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/ratelimiter/strategies"
)

type RateLimiterInterface interface {
	Check(ctx context.Context, r *http.Request) (*strategies.RateLimiterResult, error)
}

type RateLimiter struct {
	Logger              zerolog.Logger
	Strategy            strategies.LimiterStrategyInterface
	MaxRequestsPerIP    int
	MaxRequestsPerToken int
	TimeWindowMillis    int
}

func NewRateLimiter(
	logger logger.LoggerInterface,
	strategy strategies.LimiterStrategyInterface,
	ipMaxReqs int,
	tokenMaxReqs int,
	timeWindow int,
) *RateLimiter {
	return &RateLimiter{
		Logger:              logger.GetLogger(),
		Strategy:            strategy,
		MaxRequestsPerIP:    ipMaxReqs,
		MaxRequestsPerToken: tokenMaxReqs,
		TimeWindowMillis:    timeWindow,
	}
}

func (rl *RateLimiter) Check(ctx context.Context, r *http.Request) (*strategies.RateLimiterResult, error) {
	var key string
	var limit int64
	duration := time.Duration(rl.TimeWindowMillis) * time.Millisecond

	apiKey := r.Header.Get("API_KEY")

	if apiKey != "" {
		key = apiKey
		limit = int64(rl.MaxRequestsPerToken)
	} else {
		key = rip.GetClientIP(r)
		limit = int64(rl.MaxRequestsPerIP)
	}

	req := &strategies.RateLimiterRequest{
		Key:      key,
		Limit:    limit,
		Duration: duration,
	}

	result, err := rl.Strategy.Check(r.Context(), req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
