// Package main provides ...
package main

import (
	"io"
	"log"
	"net"
)

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
