package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Connections struct {
	Postgres *pgxpool.Pool
	Redis    *redis.Client
}

func Open(ctx context.Context) (*Connections, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(
		env(
			"DATABASE_URL",
			"postgres://assessment:assessment@127.0.0.1:15432/assessment?sslmode=disable",
		),
	)
	if err != nil {
		return nil, fmt.Errorf("parse database configuration: %w", err)
	}

	config.MaxConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:        env("REDIS_ADDR", "127.0.0.1:16379"),
		DialTimeout: 2 * time.Second, ReadTimeout: 2 * time.Second, WriteTimeout: 2 * time.Second,
		PoolSize: 5, ContextTimeoutEnabled: true,
	})

	connections := &Connections{Postgres: pool, Redis: client}
	if err := pool.Ping(ctx); err != nil {
		connections.Close()
		return nil, fmt.Errorf("check PostgreSQL (run make infra-up): %w", err)
	}

	if err := client.Ping(ctx).Err(); err != nil {
		connections.Close()
		return nil, fmt.Errorf("check Redis (run make infra-up): %w", err)
	}
	// Verify the supplied schema exists before a candidate starts the timer.
	if _, err := pool.Exec(ctx, "SELECT attempt_id FROM attempts LIMIT 0"); err != nil {
		connections.Close()
		return nil, fmt.Errorf("check supplied attempts table (run make migrate-up): %w", err)
	}

	return connections, nil
}

func (c *Connections) Close() {
	_ = c.Redis.Close()
	c.Postgres.Close()
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
