package redis

import (
	"context"
	"fmt"
	"time"

	redis_db "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultStartupTimeout = 30 * time.Second
	defaultPollInterval   = 100 * time.Millisecond
)

type redisWaitStrategy struct{}

var _ wait.Strategy = (*redisWaitStrategy)(nil)

func newRedisWaitStrategy() redisWaitStrategy {
	return redisWaitStrategy{}
}

func (s redisWaitStrategy) WaitUntilReady(ctx context.Context, target wait.StrategyTarget) (err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultStartupTimeout)
	defer cancel()

	if err := wait.NewHostPortStrategy(redisServicePort).WaitUntilReady(ctx, target); err != nil {
		return fmt.Errorf("error waiting for port to open: %w", err)
	}

	host, err := target.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch host: %w", err)
	}

	port, err := target.MappedPort(ctx, redisServicePort)
	if err != nil {
		return fmt.Errorf("failed to fetch port: %w", err)
	}

	return s.pollUntilReady(ctx, host, port.Int())
}

func (s redisWaitStrategy) pollUntilReady(ctx context.Context, host string, port int) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out while waiting for Redis to start: %w", ctx.Err())

		case <-time.After(defaultPollInterval):
			isReady, err := s.isReady(ctx, host, port)
			if err != nil {
				return err
			}

			if isReady {
				return nil
			}
		}
	}
}

func (s redisWaitStrategy) isReady(ctx context.Context, host string, port int) (bool, error) {
	client := redis_db.NewClient(&redis_db.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})
	defer client.Close()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return false, nil
	}

	return true, nil
}
