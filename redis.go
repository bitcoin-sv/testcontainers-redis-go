package redis

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
)

const (
	redisServicePort  = "6379/tcp"
	defaultRedisImage = "redis/redis-stack-server:6.2.6-v17"
)

type RedisContainer struct {
	testcontainers.Container
}

// RunContainer creates an instance of the Redis container type.
func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*RedisContainer, error) {
	containerRequest := testcontainers.ContainerRequest{
		Image:        defaultRedisImage,
		ExposedPorts: []string{redisServicePort},
		WaitingFor:   newRedisWaitStrategy(),
	}

	genericContainerRequest := testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerRequest); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to start Redis: %w", err)
	}

	return &RedisContainer{Container: container}, nil
}

// ServicePort returns the port on which the Redis container is listening.
func (c RedisContainer) ServicePort(ctx context.Context) (int, error) {
	port, err := c.Container.MappedPort(ctx, redisServicePort)
	if err != nil {
		return 0, err
	}

	return port.Int(), nil
}

// WithImage sets the image for the Redis container.
func WithImage(image string) testcontainers.CustomizeRequestOption {
	return testcontainers.WithImage(image)
}
