package main

import (
    "context"
    "fmt"
    "os"

    "licenseplate-plugin/internal/eventbus"
)

func main() {
    ctx := context.Background()
    fmt.Println("Starting test listener for 'events' channel...")
    err := eventbus.Listen(ctx, "events", func(ch, msg string) {
        fmt.Printf("Received on %s: %s\n", ch, msg)
        // exit if message contains 'publish_test' so demo ends
        if os.Getenv("EXIT_ON_MESSAGE") == "1" {
            os.Exit(0)
        }
    })
    if err != nil {
        fmt.Printf("Listener failed: %v\n", err)
        return
    }

    // block forever
    select {}
}
