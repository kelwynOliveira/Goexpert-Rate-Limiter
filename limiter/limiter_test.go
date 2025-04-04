package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Redis client Mock
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return redis.NewIntResult(int64(args.Int(0)), args.Error(1))
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return redis.NewBoolResult(args.Bool(0), args.Error(1))
}

// 1st scenario: Time limite expiration different for each token
func TestRateLimiter_DifferentExpirationTimes(t *testing.T) {
	mockRedis := new(MockRedisClient)
	store := NewRedisStore(mockRedis)

	token1 := "token1"
	token2 := "token2"
	limit := 5
	shortDuration := 1 * time.Second
	longDuration := 10 * time.Second

	// token1 Mock setting
	for i := 1; i <= limit; i++ {
		mockRedis.On("Incr", mock.Anything, token1).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, token1, shortDuration).Return(true, nil).Once()
		}
	}

	// token2 Mock setting
	for i := 1; i <= limit; i++ {
		mockRedis.On("Incr", mock.Anything, token2).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, token2, longDuration).Return(true, nil).Once()
		}
	}

	// token1 behavior
	for i := 1; i <= limit; i++ {
		allowed, err := store.Allow(token1, limit, shortDuration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed for token1, got %v", allowed)
	}

	// token2 behavior
	for i := 1; i <= limit; i++ {
		allowed, err := store.Allow(token2, limit, longDuration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed for token2, got %v", allowed)
	}

	mockRedis.AssertExpectations(t)
}

// Scenario 2: The access token's limit settings must override those of the IP
func TestRateLimiter_TokenOverridesIP(t *testing.T) {
	mockRedis := new(MockRedisClient)
	store := NewRedisStore(mockRedis)

	ip := "192.168.1.1"
	token := "token123"
	ipLimit := 5
	tokenLimit := 10
	duration := 1 * time.Second

	// Configure the mock for the token
	for i := 1; i <= tokenLimit; i++ {
		mockRedis.On("Incr", mock.Anything, token).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, token, duration).Return(true, nil).Once()
		}
	}

	// Configure the mock for the IP
	for i := 1; i <= ipLimit; i++ {
		mockRedis.On("Incr", mock.Anything, ip).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, ip, duration).Return(true, nil).Once()
		}
	}

	// Checks the behavior for the token
	for i := 1; i <= tokenLimit; i++ {
		allowed, err := store.Allow(token, tokenLimit, duration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed for token, got %v", allowed)
	}

	// It should block after reaching the token limit
	mockRedis.On("Incr", mock.Anything, token).Return(tokenLimit+1, nil).Once()
	allowed, err := store.Allow(token, tokenLimit, duration)
	require.NoError(t, err)
	assert.False(t, allowed, "expected blocked for token after limit, got %v", allowed)

	// Checks the behavior for the IP after configuring the token
	for i := 1; i <= ipLimit; i++ {
		allowed, err := store.Allow(ip, ipLimit, duration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed for IP, got %v", allowed)
	}

	mockRedis.AssertExpectations(t)
}

// Cenário 3: Limitação por IP
func TestRateLimiter_IPLimitation(t *testing.T) {
	mockRedis := new(MockRedisClient)
	store := NewRedisStore(mockRedis)

	ip := "192.168.1.1"
	limit := 5
	duration := 2 * time.Second

	// Configura o mock para o IP
	for i := 1; i <= limit+1; i++ {
		mockRedis.On("Incr", mock.Anything, ip).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, ip, duration).Return(true, nil).Once()
		}
	}

	// Verifica o comportamento até atingir o limite
	for i := 1; i <= limit; i++ {
		allowed, err := store.Allow(ip, limit, duration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed, got %v", allowed)
	}

	// Deve bloquear após atingir o limite
	allowed, err := store.Allow(ip, limit, duration)
	require.NoError(t, err)
	assert.False(t, allowed, "expected blocked, got %v", allowed)

	time.Sleep(duration)

	// Verifica o comportamento após o tempo de bloqueio
	mockRedis.On("Incr", mock.Anything, ip).Return(1, nil).Once()
	mockRedis.On("Expire", mock.Anything, ip, duration).Return(true, nil).Once()

	allowed, err = store.Allow(ip, limit, duration)
	require.NoError(t, err)
	assert.True(t, allowed, "expected allowed after block duration, got %v", allowed)

	mockRedis.AssertExpectations(t)
}

// Cenário 4: Limitação por token
func TestRateLimiter_TokenLimitation(t *testing.T) {
	mockRedis := new(MockRedisClient)
	store := NewRedisStore(mockRedis)

	token := "abc123"
	limit := 10
	duration := 2 * time.Second

	// Configura o mock para o token
	for i := 1; i <= limit+1; i++ {
		mockRedis.On("Incr", mock.Anything, token).Return(i, nil).Once()
		if i == 1 {
			mockRedis.On("Expire", mock.Anything, token, duration).Return(true, nil).Once()
		}
	}

	// Verifica o comportamento até atingir o limite
	for i := 1; i <= limit; i++ {
		allowed, err := store.Allow(token, limit, duration)
		require.NoError(t, err)
		assert.True(t, allowed, "expected allowed, got %v", allowed)
	}

	// Deve bloquear após atingir o limite
	allowed, err := store.Allow(token, limit, duration)
	require.NoError(t, err)
	assert.False(t, allowed, "expected blocked, got %v", allowed)

	time.Sleep(duration)

	// Verifica o comportamento após o tempo de bloqueio
	mockRedis.On("Incr", mock.Anything, token).Return(1, nil).Once()
	mockRedis.On("Expire", mock.Anything, token, duration).Return(true, nil).Once()

	allowed, err = store.Allow(token, limit, duration)
	require.NoError(t, err)
	assert.True(t, allowed, "expected allowed after block duration, got %v", allowed)

	mockRedis.AssertExpectations(t)
}
