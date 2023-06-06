package main

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
  time.Sleep(10 * time.Second)
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	// Access JetStream
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error accessing JetStream: %v", err)
	}

	// Setup a durable, concurrent subscriber on the "SYSLOGS.sources" subject.
	// The "nats" argument specifies the queue group.
	var count int64
	sub, err := js.QueueSubscribe("SYSLOGS.sources", "nats", func(msg *nats.Msg) {
		atomic.AddInt64(&count, 1)
	})
	if err != nil {
		log.Fatalf("Error subscribing to SYSLOGS.sources: %v", err)
	}
	defer sub.Unsubscribe() // Unsubscribe when done.

	// Setup a ticker to print the number of messages received every second.
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			val := atomic.SwapInt64(&count, 0)
			log.Printf("Processed %d messages in the last second", val)
		}
	}()

	// Setup a channel to handle OS signals to gracefully shut down
	// when an interrupt signal is received.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// Cleanup when done. Unsubscribe to the subject.
	sub.Unsubscribe()
	nc.Drain()

	log.Println("Connection to NATS closed.")
}

