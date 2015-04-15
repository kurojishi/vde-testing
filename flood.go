package main

import (
	"flag"
	"fmt"
	"net"
)

func sendData(adr string) {
	var payload = []byte("skldjaslkcjlak")
	conn, err := net.Dial("tcp", adr)
	if err != nil {
		fmt.Print(err)
	}
	_, err = conn.Write(payload)

	fmt.Printf("Sent: ", payload)
	conn.Close()
}

func receiveData(conn net.Listener, ok chan bool) {
	fmt.Println("receiver starterd")
	payload := make([]byte, 14)
	listenerConnection, err := conn.Accept()
	if err != nil {
		fmt.Print(err)
	}
	a, err := listenerConnection.Read(payload)
	fmt.Print(a)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Received: ", payload)
	conn.Close()
	ok <- true

}

func main() {
	var in = flag.String("in", "192.168.4.1:5000", "address to send the data too")
	//var out = flag.String("out", "192.168.4.2:5000", "address to send the data from")
	//outAddress := net.ParseIP(*out)
	//if outAddress == nil {
	//log.Fatalln("bad output ip address")
	//}
	listener, err := net.Listen("tcp", *in)
	defer listener.Close()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("server started")
	fmt.Println("connection accepted")
	//inAddress := net.ParseIP(*in)
	//if inAddress == nil {
	//log.Fatalln("bad input ip address")
	//}
	ok := make(chan bool)
	go receiveData(listener, ok)
	go sendData(*in)
	x := <-ok
	if !x {
		fmt.Println("WTF")
	}
}
