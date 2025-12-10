package eventbus

import (
    "context"
    "os"

    "github.com/redis/go-redis/v9"
)

func newClient() *redis.Client {
    addr := os.Getenv("HUB_BUS_ADDR")
    if addr == "" {
        addr = "hub_bus:6379"
    }

    return redis.NewClient(&redis.Options{
        Addr: addr,
        DB:   0,
    })
}

// Publish publishes a string message to the given channel on the event bus.
func Publish(ctx context.Context, channel string, message string) error {
    r := newClient()
    return r.Publish(ctx, channel, message).Err()
}
