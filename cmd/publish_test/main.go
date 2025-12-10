package main

import (
    "context"
    "fmt"
    "time"

    "licenseplate-plugin/internal/eventbus"
)

func main() {
    fmt.Println("Publishing test message to 'events' channel...")
    ctx := context.Background()
    err := eventbus.Publish(ctx, "events", "{\"test\":\"hello from publish_test\"}")
    if err != nil {
        fmt.Printf("Publish failed: %v\n", err)
        return
    }
    fmt.Println("Published. Waiting a bit before exit...")
    time.Sleep(1 * time.Second)
}
