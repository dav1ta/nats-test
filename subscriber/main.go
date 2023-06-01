package main

import (
    "log"
    "strings"
    "github.com/nats-io/nats.go"
    "time"
)

func extractIP(input string) string {
    split := strings.Split(input, "srcIP=")
    if len(split) < 2 {
        return ""
    }
    return strings.Split(split[1], "|")[0]
}

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

    _, err = js.Subscribe("SYSLOGS.sources", func(m *nats.Msg) {
        log.Printf("Received a message: %s\n", string(m.Data))
        ip := extractIP(string(m.Data))
        if ip != "" {
            log.Printf("Extracted IP: %s\n", ip)
        } else {
            log.Printf("No IP found in the message.")
        }
    }, nats.Durable("my-durable"))

    if err != nil {
        log.Fatalf("Error subscribing to ORDERS.order1: %v", err)
    }

    select {}
}

