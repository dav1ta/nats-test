package main

import (
    "fmt"
    "log"
    "math/rand"
    "time"

    "github.com/nats-io/nats.go"
)

const (
    Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    Digits  = "0123456789"
)

func RandStringBytes(n int, charset string) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

func RandomIP() string {
    return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func GenerateNewSyslog() string {
    timestamp := time.Now().Format("Mon Jan _2 15:04:05 2006")
    hostname := RandStringBytes(5, Letters)
    firewall := RandStringBytes(5, Letters)
    srcIP := RandomIP()
    dstIP := RandomIP()

    return fmt.Sprintf("%s: <14>Mar  4 15:53:03 %s %s/box_Firewall_Activity:  Info     %s Remove: type=FWD|proto=UDP|srcIF=eth1|srcIP=%s|srcPort=35119|srcMAC=08:00:27:da:d7:9c|dstIP=%s|dstPort=53|dstService=domain|dstIF=eth0|rule=InternetAccess/<App>:RestrictTim|info=Balanced Session Idle Timeout|srcNAT=%s|dstNAT=%s|duration=21132|count=1|receivedBytes=130|sentBytes=62|receivedPackets=1|sentPackets=1|user=|protocol=|application=|target=|content=|urlcat",
        timestamp, hostname, firewall, firewall, srcIP, dstIP, srcIP, dstIP)
}

func main() {
    rand.Seed(time.Now().UnixNano())

    nc, _ := nats.Connect("nats://nats:4222")
    js, _ := nc.JetStream()

  for i := 0; i < 1; i++ {
    a:=0

    for i := 0; i < 500000; i++ {
        a++
        msg := GenerateNewSyslog()
        // log.Printf("Publishing message: %s to stream 'SYSLOGS'\n", msg)
      log.Println(a)

        _, err := js.Publish("SYSLOGS.sources", []byte(msg))
        if err != nil {
            log.Fatal(err)
        }
    }
    time.Sleep(1 * time.Second)

  }

    nc.Drain()
}

