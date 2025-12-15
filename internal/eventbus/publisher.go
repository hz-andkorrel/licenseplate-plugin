package eventbus

import (
    "context"
    "log"

    "github.com/redis/go-redis/v9"
)

// Publish publishes a string message to the given channel on the event bus.
// Uses the provided Redis client (typically a global shared client).
func Publish(ctx context.Context, client *redis.Client, channel string, message string) error {
    if client == nil {
        log.Println("Redis client is nil, skipping event publish")
        return nil
    }
    return client.Publish(ctx, channel, message).Err()
}
