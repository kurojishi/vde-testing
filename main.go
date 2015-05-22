package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

func main() {
	var server bool
	var address, iface string
	var port int
	var size, snaplen int64
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "p", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.Int64Var(&snaplen, "s", 1600, "spanlen for pcap capture")
	flag.Int64Var(&size, "size", 150, "ho much data to send")
	flag.Parse()
	if server {
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		ok := make(chan bool)
		ready := make(chan int)
		go receiveData("tcp", address+":"+strconv.Itoa(port), ready)
		go StreamStats(iface, snaplen, strconv.Itoa(port), ready)
		x := <-ok
		if !x {
			log.Fatalf("WTF")
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
