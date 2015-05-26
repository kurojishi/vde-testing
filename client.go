package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const (
	kb int64 = 1000
	mb int64 = 1000 * kb
	gb int64 = 1000 * mb
)

//controlServer start the controls channel on the client
func controlServer(bind, address string) {
	clistener, err := net.Listen("tcp", bind)
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
			sendData(address, 1000)
			//case load:
			//case stress:
		case die:
			break

		default:
			continue
		}
	}
}

//sendData send size data (in megabytes)to the string addr
func sendData(addr string, size int64) {
	log.Println("sending data")
	_, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("sendData: %v", err)
	}
	n, err := io.CopyN(conn, devZero, size*(mb))
	if err != nil {
		log.Fatal(err)
	}
	if n != size*mb {
		log.Fatalf("couldnt send %v Megabytes", float64(n)/float64(mb))
	}
	log.Printf("sent %v MB", size)
}
