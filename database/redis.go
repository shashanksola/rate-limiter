package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shashanksola/rate-limiter/config"
)

// NewRedisClient creates a new Redis client and verifies the connection.

func NewRedisClient(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		// Pool configuration (similar to database/sql pool tuning)
		PoolSize:     20,              // Max connections in the pool
		MinIdleConns: 5,               // Keep at least 5 idle connections warm
		MaxRetries:   3,               // Retry failed commands up to 3 times
		DialTimeout:  5 * time.Second, // Timeout for establishing new connections
		ReadTimeout:  3 * time.Second, // Timeout for reading a response
		WriteTimeout: 3 * time.Second, // Timeout for writing a command
	})

	// Verify the connection with a PING command.
	// Redis PING returns "PONG" if the server is alive.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pingRedisWithRetry(ctx, client); err != nil {
		return nil, fmt.Errorf("connecting to redis: %w", err)
	}

	slog.Info("connected to Redis",
		"addr", cfg.Addr,
		"db", cfg.DB,
	)

	return client, nil
}

// pingRedisWithRetry attempts to ping Redis with retry logic.
// Same pattern as the PostgreSQL retry — production systems need resilience.
func pingRedisWithRetry(ctx context.Context, client *redis.Client) error {
	var lastErr error
	maxRetries := 5
	backoff := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		if err := client.Ping(ctx).Err(); err != nil {
			lastErr = err
			slog.Warn("redis ping failed, retrying...",
				"attempt", i+1,
				"max_retries", maxRetries,
				"error", err,
			)

			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled while waiting to retry: %w", ctx.Err())
			case <-time.After(backoff):
				backoff *= 2
			}
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to ping redis after %d attempts: %w", maxRetries, lastErr)
}
