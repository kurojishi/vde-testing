package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const (
	ready     int32 = 0
	bandwidth int32 = 1
	latency   int32 = 2
	load      int32 = 3
	stress    int32 = 4
)

type zeroFile struct{}

func (d *zeroFile) Read(p []byte) (int, error) {
	return len(p), nil
}

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
	listener, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	log.Println("server started")
	synch <- ready
	cch <- bandwidth
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
