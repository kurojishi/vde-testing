package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/kurojishi/govde"
)

func sendData(conn *govde.ConnectionVde) {
	for true {
		var payload = []byte("skldjaslkcjlak")
		conn.Send(payload)

		fmt.Println("Sent: ", payload)
	}
}

func receiveData(conn *govde.ConnectionVde) {
	for true {
		payload, err := conn.Receive()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Received: ", payload)
	}

}

func main() {
	var in = flag.String("in", "", "address to send the data too")
	var out = flag.String("out", "", "address to send the data from")
	var configuration map[string]string
	outAddress := net.ParseIP(*out)
	if outAddress == nil {
		log.Fatalln("bad output ip address")
	}
	OutConnection, err := govde.Connect("TCP", outAddress.String(), configuration)
	if err != nil {
		log.Fatal(err)
	}
	inAddress := net.ParseIP(*in)
	if inAddress == nil {
		log.Fatalln("bad input ip address")
	}
	listener, err := govde.Listen("interface", inAddress.String(), configuration)
	if err != nil {
		log.Fatal(err)
	}
	go receiveData(listener)
	go sendData(OutConnection)
}
