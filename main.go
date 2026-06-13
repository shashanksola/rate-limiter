package ratelimiter

import (
	"log/slog"
	"os"

	"github.com/shashanksola/rate-limiter/config"
	"github.com/shashanksola/rate-limiter/database"
)

func main() {

	// Initialize logger

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting GOREST API server")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("configuration loaded",
		"server_port", cfg.Server.Port,
		"redis_addr", cfg.Redis.Addr,
	)

	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		slog.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing Redis connection")
		redisClient.Close()
	}()
}
