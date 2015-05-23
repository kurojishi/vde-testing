// Package main provides ...
package main

import (
	"bytes"
	"log"
	"net"
)

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
		if kind == "client" {
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
		} else if kind == "server" {

		} else {
			log.Fatal("Wrong kind of controlServer")
		}
	}
}
