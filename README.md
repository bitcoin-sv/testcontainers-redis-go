<!-- markdownlint-configure-file {
  "MD033": false,
  "MD041": false
} -->

<div align="center">

# testcontainers-redis-go

[![Go Reference](https://pkg.go.dev/badge/github.com/bitcoin-sv/testcontainers-redis-go.svg)](https://pkg.go.dev/github.com/bitcoin-sv/testcontainers-redis-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitcoin-sv/testcontainers-redis-go)](https://goreportcard.com/report/github.com/bitcoin-sv/testcontainers-redis-go)

Go library for **[Redis](https://redis.io/) integration testing via
[Testcontainers](https://testcontainers.com/)**.

</div>

## Install

Use `go get` to install the latest version of the library.

```bash
go get -u github.com/bitcoin-sv/testcontainers-redis-go@latest
```

## Usage

```go
import (
    "context"
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
    redis_db "github.com/redis/go-redis/v9"
    redisTest "github.com/bitcoin-sv/testcontainers-redis-go"
)

func TestRedis(t *testing.T) {
    redisClient := setupRedis(t)
    // your code here
}

func setupRedis(t *testing.T) *redis_db.Client {
    ctx := context.Background()

    container, err := redisTest.RunContainer(ctx)
    require.NoError(t, err)
    
    t.Cleanup(func() {
        err := container.Terminate(ctx)
        require.NoError(t, err)
    })

    host, err := container.Host(ctx)
    require.NoError(t, err)
    
    port, err := container.ServicePort(ctx)
    require.NoError(t, err)

    client, err := redis_db.NewClient(&redis_db.Options{
        Addr: fmt.Sprintf("%s:%d", host, port),
    })
    require.NoError(t, err)

    return client
}
```
