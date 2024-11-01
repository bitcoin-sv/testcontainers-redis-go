package redis

import (
	"context"
	"fmt"
	"testing"

	redis_db "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	ctx := context.Background()

	container, err := RunContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoErrorf(t, err, "failed to terminate Redis container")
	})

	host, err := container.Host(ctx)
	require.NoErrorf(t, err, "failed to fetch Redis host")
	port, err := container.ServicePort(ctx)
	require.NoErrorf(t, err, "failed to fetch Redis port")

	client := redis_db.NewClient(&redis_db.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
	})
	defer client.Close()

	err = client.Set(ctx, "key", "value", 0).Err()
	require.NoErrorf(t, err, "failed to set Redis key")

	val, err := client.Get(ctx, "key").Result()
	require.NoErrorf(t, err, "failed to get Redis key")
	require.Equal(t, "value", val)
}
