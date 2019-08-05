package main

/*
Program:    main.go
Component:  mqtttest
Language:   go
Support:    David A. Fahey - whome God preserve.
Purpose:    To test a Go based MQTT.

History:    04Aug2019 Initial coding                                    DAF

Building:
	        cd /home/david/usr/GoApp/gopath/src/github.com/x0ray
            git clone 'https://github.com/x0ray/tmq'
            cd tmq
	        go mod init github.com/x0ray/tmq
            go build ./...

            // non-module aware alternate
	        // GO111MODULE=off;go install ./...

Testing:
            /home/david/gopath/bin/tmq

Notes:      This code is based on the example found here:
                https://github.com/DrmagicE/gmqtt

Output:


*/

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/pkg/packets"
)

const (
	PGM     = "tmq.go"
	VER     = "0.0.1"
	VERDATE = "04Aug2019"
)

func main() {
	log.Printf("INFO Server: %s ver: %s of: %s at: %v\n", PGM, VER, VERDATE, time.Now())

	s := gmqtt.NewServer()

	ln, err := net.Listen("tcp", ":1883")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	crt, err := tls.LoadX509KeyPair("../testcerts/server.crt", "../testcerts/server.key")
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsln, err := tls.Listen("tcp", ":8883", tlsConfig)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	s.AddTCPListenner(ln)
	s.AddTCPListenner(tlsln)
	//Configures and registers callback before s.Run()
	s.SetMaxInflightMessages(20)
	s.SetMaxQueueMessages(99999)
	s.RegisterOnSubscribe(func(client *gmqtt.Client, topic packets.Topic) uint8 {
		if topic.Name == "test/nosubscribe" {
			return packets.SUBSCRIBE_FAILURE
		}
		return topic.Qos
	})
	s.Run()
	log.Printf("INFO Server: %s started.\n", PGM)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	s.Stop(context.Background())

	log.Printf("INFO Server: %s ended.\n", PGM)
}
