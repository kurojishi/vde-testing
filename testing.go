package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"time"
)

func sendData(addr string, ok chan bool) {
	log.Println("sending data")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("sendData: %v", err)
	}
	time.Sleep(time.Duration(100000))
	payload := make([]byte, 40960)
	payload = []byte("AAAAAAAAAAAAARGH")
	_, err = conn.Write(payload)
	ok <- true
}

func receiveData(conn net.Listener) {
	log.Println("receiver starterd")
	for {
		listenerConnection, err := conn.Accept()
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, listenerConnection)
		//_, err := listenerConnection.Read(payload)
		if err != nil {
			log.Fatalf("receiveData: %v", err)
		}
		log.Printf("data received: %v", buf.String())
	}

}

func printStats(stream chan StatsStream) {
	for s := range stream {
		diffSecs := float64(s.end.Sub(s.start)) / float64(time.Second)
		log.Printf("Reassembly of stream %v:%v complete - start:%v end:%v bytes:%v packets:%v ooo:%v bps:%v pps:%v skipped:%v",
			s.net, s.transport, s.start, s.end, s.bytes, s.packets, s.outOfOrder,
			float64(s.bytes)/diffSecs, float64(s.packets)/diffSecs, s.skipped)

	}
}

func main() {
	var server bool
	var address, port, iface string
	var snaplen int
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.StringVar(&port, "p", "5000", "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.IntVar(&snaplen, "s", 1600, "spanlen for pcap capture")
	flag.Parse()
	if server {
		defer close(statsResults)
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		listener, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			log.Fatalf("main; %v", err)
		}
		log.Println("server started")
		ok := make(chan bool)
		go StreamStats(iface, int32(snaplen))
		go printStats(statsResults)
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
		go sendData(address+":"+port, ok)
		x := <-ok
		if !x {
			log.Println("WTF")
		}
	}
}
