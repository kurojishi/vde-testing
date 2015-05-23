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
	var snaplen int64
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "p", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.Int64Var(&snaplen, "snaplen", 1600, "spanlen for pcap capture")
	flag.Parse()
	fullAddr := address + ":" + strconv.Itoa(port)
	if server {
		sPort := strconv.Itoa(port)
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		ready := make(chan int)
		go StreamStats(iface, snaplen, sPort, ready)
		receiveData("tcp", fullAddr, ready)
	} else {
		sendData(fullAddr, 150)
	}
}
