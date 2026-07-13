package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ai-finops/ai-finops/packages/shared"
	"github.com/redis/go-redis/v9"
)

// RateLimiter enforces per-API-key ingest limits using Redis.
type RateLimiter struct {
	client *redis.Client
	limit  int
	log    *slog.Logger
}

// NewRateLimiter connects to Redis when configured; otherwise it no-ops.
func NewRateLimiter(redisURL string, limit int, log *slog.Logger) (*RateLimiter, error) {
	limiter := &RateLimiter{limit: limit, log: log}
	if redisURL == "" {
		log.Warn("redis url not configured; ingestion rate limiting disabled")
		return limiter, nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Warn("redis unavailable; ingestion rate limiting disabled", slog.String("error", err.Error()))
		_ = client.Close()
		return limiter, nil
	}

	limiter.client = client
	return limiter, nil
}

// Allow returns whether the API key is within the configured per-minute limit.
func (r *RateLimiter) Allow(ctx context.Context, keyID string) error {
	if r.client == nil || keyID == "" {
		return nil
	}

	bucket := time.Now().UTC().Format("200601021504")
	redisKey := fmt.Sprintf("ingestion:ratelimit:%s:%s", keyID, bucket)

	count, err := r.client.Incr(ctx, redisKey).Result()
	if err != nil {
		r.log.Warn("rate limit check failed; allowing request", slog.String("error", err.Error()))
		return nil
	}
	if count == 1 {
		_ = r.client.Expire(ctx, redisKey, time.Minute).Err()
	}
	if count > int64(r.limit) {
		return shared.RateLimited("rate_limited", "api key exceeded ingest rate limit")
	}
	return nil
}

// Close releases the Redis client.
func (r *RateLimiter) Close() error {
	if r.client == nil {
		return nil
	}
	return r.client.Close()
}