package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
  "time"
	"github.com/nats-io/nats.go"
)




type SyslogEvent struct {
	EventID         string
	EventStatus     string
	DeviceName      string
	Protocol        string
	SrcIP           string
	DestIP          string
	Rule            string
	Duration        int
	Count           int
	ReceivedBytes   int
	SentBytes       int
	ReceivedPackets int
}

type Rule struct {
    Name        string
    Description string
    Condition   func(e SyslogEvent) bool
}
func main() {
	// Connect to NATS

var deviceFailures = make(map[string]int)

    rules := []Rule{
      {
        Name:        "Potential Brute Force Attack",
        Description: "More than 5 login failures on a device within a minute",
        Condition: func(e SyslogEvent) bool {
          if e.EventID == "LOGIN" && e.EventStatus == "FAILURE" {
            deviceFailures[e.DeviceName]++
            if deviceFailures[e.DeviceName] > 2 {
              // Reset the failure count after triggering the alert
              deviceFailures[e.DeviceName] = 0
              return true
            }
          }
          return false
        },
      },
    }


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
	sub, err := js.Subscribe("SYSLOGS.sources", func(msg *nats.Msg) {
		re := regexp.MustCompile(`.* (\S+)\/.*type=(\S+)\|.*srcIP=(\S+)\|.*dstIP=(\S+)\|.*rule=(\S+)\|.*duration=(\d+)\|.*count=(\d+)\|.*receivedBytes=(\d+)\|.*sentBytes=(\d+)\|.*receivedPackets=(\d+).*`)
		match := re.FindStringSubmatch(string(msg.Data))

		if len(match) == 0 {
			log.Println("No match found in message")
			return
		}

		event := SyslogEvent{}
		event.DeviceName = match[1]
		event.Protocol = match[2]
		event.SrcIP = match[3]
		event.DestIP = match[4]
		event.Rule = match[5]
		event.Duration, _ = strconv.Atoi(match[6])
		event.Count, _ = strconv.Atoi(match[7])
		event.ReceivedBytes, _ = strconv.Atoi(match[8])
		event.SentBytes, _ = strconv.Atoi(match[9])
		event.ReceivedPackets, _ = strconv.Atoi(match[10])

		log.Printf("Received event: %+v", event)
    for _, rule := range rules {
        if rule.Condition(event) {
            log.Printf("Alert! Rule '%s' triggered: %s", rule.Name, rule.Description)
        }
}
	})

	if err != nil {
		log.Fatalf("Error subscribing to SYSLOGS.sources: %v", err)
	}

	// Setup a channel to handle OS signals to gracefully shut down
	// when an interrupt signal is received.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	// Cleanup when done. Unsubscribe to the subject.
	sub.Unsubscribe()
	nc.Drain()


  go func() {
    for range time.Tick(1 * time.Minute) {
      for device := range deviceFailures {
        deviceFailures[device] = 0
      }
    }
  }()

	log.Println("Connection to NATS closed.")
}

