package eventbus

import (
    "context"
    "fmt"
    "os"

    "github.com/redis/go-redis/v9"
)

func Listen(ctx context.Context, channel string, handler func(channel, message string)) error {
    addr := os.Getenv("HUB_BUS_ADDR")
    if addr == "" {
        addr = "hub_bus:6379"
    }

    r := redis.NewClient(&redis.Options{
        Addr: addr,
        DB:   0,
    })

    sub := r.Subscribe(ctx, channel)
    ch := sub.Channel()

    go func() {
        for msg := range ch {
            handler(msg.Channel, msg.Payload)
        }
    }()

    // Return nil immediately; caller should manage context cancellation
    fmt.Printf("Subscribed to channel: %s\n", channel)
    return nil
}
