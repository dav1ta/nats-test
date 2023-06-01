package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	StreamName     = "SYSLOGS"
	StreamSubjects = "SYSLOGS.*"
)

func CreateStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(StreamName)

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	err = CreateStream(js)
	if err != nil {
		log.Fatal(err)
	}

//use StreamName and StreamSubjects to create log
log.Println("Stream: ",StreamName,"created")
  

}

