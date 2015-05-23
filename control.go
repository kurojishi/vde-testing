package main

import (
	"bytes"
	"log"
	"net"
)

//controlServer start the controls channel on the client
func controlServer(address string, cch chan int) {
	clistener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Control Server Started")
	for {
		conn, err := clistener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		var buf bytes.Buffer
		buf.ReadFrom(conn)
		//TODO: define the other cases
		switch buf.String() {
		case "ready":
			cch <- 1
		case "bandwidth":
			cch <- 2
		case "latency":
			cch <- 3
		case "load":
			cch <- 4
		case "stress":
			cch <- 5
		}
	}
}

func sendControlSignal(address string, msg string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}
