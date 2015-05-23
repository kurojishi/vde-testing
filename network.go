package main

import (
	"bytes"
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

//controlServer start the controls channel from the client to the server and vice versa
func controlServer(kind string, address string, cch chan int) {
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
