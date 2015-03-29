package main

import (
	"flag"
	"github.com/kurojishi/govde"
	"log"
	"net"
)

func sendData(conn *govde.ConnectionVde) {
	for true {
		var payload = []byte("skldjaslkcjlak")
		conn.Send(payload)

	}
}

func receiveData(conn *govde.ConnectionVde) {
	for true {
		_, err := conn.Receive()
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	var in = flag.String("in", "", "interface that will receive data")
	var out = flag.String("out", "", "interface that will send data")
	var configuration map[string]string
	outFace, err := net.InterfaceByName(*out)
	if err != nil {
		log.Fatal(err)
	}
	inFace, err := net.InterfaceByName(*in)
	if err != nil {
		log.Fatal(err)
	}
	outAddress, err := outFace.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	OutConnection, err := govde.Connect("TCP", outAddress[1].String(), configuration)
	if err != nil {
		log.Fatal(err)
	}
	inAddress, err := inFace.Addrs()
	if err != nil {
		log.Fatal(err)
	}
	listener, err := govde.Listen("interface", inAddress[1].String(), configuration)
	if err != nil {
		log.Fatal(err)
	}
	go receiveData(listener)
	go sendData(OutConnection)
}
