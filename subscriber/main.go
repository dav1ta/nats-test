package main

import (
    "log"
    "github.com/nats-io/nats.go"
  "time"
)

func main() {
    nc, err := nats.Connect("nats://nats:4222")
    if err != nil {
        log.Fatalf("Error connecting to NATS: %v", err)
    }
  time.Sleep(5 * time.Second)
    
    js, err := nc.JetStream()
    if err != nil {
        log.Fatalf("Error accessing JetStream: %v", err)
    }

    // Create a JetStream subscription
    _, err = js.Subscribe("ORDERS.order1", func(m *nats.Msg) {
        // Process each message
        log.Printf("Received a message: %s\n", string(m.Data))
    }, nats.Durable("my-durable"))
    if err != nil {
        log.Fatalf("Error subscribing to ORDERS.order1: %v", err)
    }

    // Keep the connection alive until the program is terminated
    select {}
}

