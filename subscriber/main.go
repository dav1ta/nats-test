package main

import (
    "log"
    "github.com/nats-io/nats.go"
    "time"
    "sync/atomic"
)

func main() {
    nc, err := nats.Connect("nats://nats:4222")
    if err != nil {
        log.Fatalf("Error connecting to NATS: %v", err)
    }
    time.Sleep(50 * time.Second)

    js, err := nc.JetStream()
    if err != nil {
        log.Fatalf("Error accessing JetStream: %v", err)
    }

    sub, err := js.PullSubscribe("SYSLOGS.sources", "my-durable")
    if err != nil {
        log.Fatalf("Error subscribing to SYSLOGS.sources: %v", err)
    }

    var count int64

    go func() {
        for {
            time.Sleep(time.Second)
            val := atomic.SwapInt64(&count, 0)
            log.Printf("Processed %d messages in the last second", val)
        }
    }()

    for {
        msgs, err := sub.Fetch(100, nats.MaxWait(1*time.Second))
        if err != nil {
            log.Printf("Error fetching messages: %v", err)
            continue
        }

        go func(msgs []*nats.Msg) {
            for range msgs {
                atomic.AddInt64(&count, 1)
            }
        }(msgs)
    }
}
