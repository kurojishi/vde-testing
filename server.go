package main

import (
	"io"
	"log"
	"net"
)

type nullFile struct{}
type zeroFile struct{}

func (d *nullFile) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d *zeroFile) Read(p []byte) (int, error) {
	return len(p), nil
}

var devZero = &zeroFile{}
var devNull = &nullFile{}

func receiveData(protocol string, address string, cch chan int) {
	if ready := <-cch; ready != 1 {
		log.Fatal("Wrong control message")
	}
	listener, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	log.Println("server started")
	log.Println("receiver starterd")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
		_, err = io.Copy(conn, devZero)
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
	}

}
