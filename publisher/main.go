
package main

import (
    "log"
    "github.com/nats-io/nats.go"
  "time"
  "fmt"
)

func main() {
    nc, _ := nats.Connect("nats://nats:4222")
    js, _ := nc.JetStream()

    for i := 0; i < 10; i++ {
        msg := fmt.Sprintf("Order %d", i)
      log.Printf("Publishing message: %s to stream 'ORDERS.order1'\n", msg)
        _, err := js.Publish("ORDERS.order1", []byte(msg))
        time.Sleep(1 * time.Second)
        if err != nil {
            log.Fatal(err)
        }
    }

    nc.Drain()
} 

