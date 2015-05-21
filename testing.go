package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strconv"
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

//sendData send size data (in megabytes)to the string addr
func sendData(addr string, size int64, ok chan bool) {
	log.Println("sending data")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("sendData: %v", err)
	}
	n, err := io.CopyN(conn, devZero, size*(1000000))
	if err != nil {
		log.Fatal(err)
	}
	if n != size*1000000 {
		log.Fatalf("couldnt send %v Megabytes", float64(n)/float64(1000000))
	}
	log.Printf("sent %v MB", size)
	ok <- true
}

func receiveData(conn net.Listener) {
	log.Println("receiver starterd")
	for {
		listenerConnection, err := conn.Accept()
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
		_, err = io.Copy(devNull, listenerConnection)
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
	}

}

func main() {
	var server bool
	var address, iface string
	var port, snaplen int
	var size int64
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "p", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.IntVar(&snaplen, "s", 1600, "spanlen for pcap capture")
	flag.Int64Var(&size, "size", 150, "ho much data to send")
	flag.Parse()
	if server {
		defer close(statsResults)
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		listener, err := net.Listen("tcp", address+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("main; %v", err)
		}
		log.Println("server started")
		ok := make(chan bool)
		go StreamStats(iface, int32(snaplen), port)
		if err != nil {
			log.Fatalf("main; %v", err)
		}
		go receiveData(listener)
		x := <-ok
		if !x {
			log.Fatalf("main; %v", err)
		}
	} else {
		ok := make(chan bool)
		go sendData(address+":"+strconv.Itoa(port), size, ok)
		x := <-ok
		if !x {
			log.Println("WTF")
		}
	}
}
