package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const (
	bandwidth int32 = 1
	latency   int32 = 2
	load      int32 = 3
	stress    int32 = 4
)

const (
	stop  int32 = 0
	ready int32 = 1
)

type zeroFile struct{}

type nullFile struct{}

func (d *nullFile) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d *zeroFile) Read(p []byte) (int, error) {
	return len(p), nil
}

var devNull = &nullFile{}
var devZero = &zeroFile{}

func sendControlSignal(address string, msg int32) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	err = binary.Write(conn, binary.LittleEndian, msg)
	if err != nil {
		return err
	}
	log.Print("send control message")
	return nil
}

func signalLoop(control string, cch chan int32) {
	for msg := range cch {
		err := sendControlSignal(control, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func receiveData(protocol string, address string, cch, synch chan int32) {
	<-synch
	listener, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	cch <- bandwidth
	log.Println("server started, control message sent")
	for {
		conn, err := listener.Accept()
		log.Print("accepted connection")
		if err != nil {
			log.Fatalf("connection error: %v", err)
		}
		_, err = io.Copy(devNull, conn)
		if err != nil {
			log.Fatalf("data receive error: %v", err)
		}
		<-synch
	}

}
