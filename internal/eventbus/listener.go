package eventbus

import (
    "context"
    "fmt"
    "log"

    "github.com/redis/go-redis/v9"
)

// Listen subscribes to the given channel and calls the handler for each message.
// Uses the provided Redis client (typically a global shared client).
func Listen(ctx context.Context, client *redis.Client, channel string, handler func(channel, message string)) error {
    if client == nil {
        log.Println("Redis client is nil, skipping event listener")
        return nil
    }

    sub := client.Subscribe(ctx, channel)
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
