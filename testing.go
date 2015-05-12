package main

import (
	"flag"
	"log"
	"net"
)

func sendData(addr string, ok chan bool) {
	payload := make([]byte, 4096)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sending data")
	for i := 0; i < 20; i++ {
		_, err = conn.Write(payload)
	}
	conn.Close()
	ok <- true
}

func receiveData(conn net.Listener) {
	log.Println("receiver starterd")
	payload := make([]byte, 4096)
	listenerConnection, err := conn.Accept()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		_, err := listenerConnection.Read(payload)
		log.Println(string(payload))
		if err != nil {
			log.Fatal(err)
		}
	}
	conn.Close()

}

func main() {
	var server bool
	var address string
	var port string
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.StringVar(&port, "p", "5000", "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.Parse()
	if server {
		listener, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("server started")
		ok := make(chan bool)
		if err != nil {
			log.Fatal(err)
		}
		go receiveData(listener)
		x := <-ok
		if !x {
			log.Fatal(err)
		}
	} else {
		ok := make(chan bool)
		go sendData(address+":"+port, ok)
		x := <-ok
		if !x {
			log.Println("WTF")
		}
	}
}
