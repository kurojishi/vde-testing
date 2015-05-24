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
	die       int32 = 0
)

const (
	stop  int32 = 1
	ready int32 = 2
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

func bandwidthTest(iface, port, address string, snaplen int64, cch chan int32) {
	sync := make(chan int32)
	go TCPStats(iface, snaplen, port, sync)
	<-sync
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	cch <- bandwidth
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	_, err = io.Copy(devNull, conn)
	if err != nil {
		log.Fatalf("data receive error: %v", err)
	}
	<-sync

}
