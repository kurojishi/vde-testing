// Package main provides ...
package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type nullFile struct{}

func (d *nullFile) Write(p []byte) (int, error) {
	return len(p), nil
}

var devNull = &nullFile{}

//controlServer start the controls channel on the client
func controlServer(address string, cch chan int32) {
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
		var buf int32
		binary.Read(conn, binary.LittleEndian, &buf)
		log.Printf("control message arrived")
		//TODO: define the other cases
		switch buf {
		case bandwidth:
			cch <- 1
		case latency:
			cch <- 2
		case load:
			cch <- 3
		case stress:
			cch <- 4
		default:
			continue
		}
	}
}

func startTests(cch chan int32, addr string) {
	for msg := range cch {
		switch msg {
		case bandwidth:
			log.Println("starting bandwidth test")
			sendData(addr, 150)
		default:
			continue

		}
	}

}

//sendData send size data (in megabytes)to the string addr
func sendData(addr string, size int64) {
	log.Println("sending data")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("sendData: %v", err)
	}
	n, err := io.CopyN(devNull, conn, size*(1000000))
	if err != nil {
		log.Fatal(err)
	}
	if n != size*1000000 {
		log.Fatalf("couldnt send %v Megabytes", float64(n)/float64(1000000))
	}
	log.Printf("sent %v MB", size)
}
