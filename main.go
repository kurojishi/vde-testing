package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

func main() {
	var server bool
	var port int
	var snaplen int64
	var address, control, iface string
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "p", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.Int64Var(&snaplen, "snaplen", 1600, "spanlen for pcap capture")
	flag.StringVar(&control, "b", "192.168.4.15:8000", "")
	flag.Parse()
	fullAddr := address + ":" + strconv.Itoa(port)
	if server {
		sPort := strconv.Itoa(port)
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		sync := make(chan int32)
		cch := make(chan int32)
		go signalLoop(control, cch)
		go StreamStats(iface, snaplen, sPort, sync)
		receiveData("tcp", fullAddr, cch, sync)
	} else {
		cch := make(chan int32)
		controlServer(control, cch)
		go startTests(cch, fullAddr)
	}
}
