package main

import (
	"flag"
	"log"
	"net"
	"strconv"
)

var pid int

func main() {
	var server bool
	var port int
	var snaplen int64
	var address, remote, iface string
	flag.StringVar(&address, "addr", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "port", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.Int64Var(&snaplen, "snaplen", 1600, "spanlen for pcap capture")
	flag.StringVar(&remote, "raddr", "192.168.4.15", "")
	flag.IntVar(&pid, "pid", 0, "the vde switch pid")
	flag.Parse()
	fullAddr := address + ":" + strconv.Itoa(port)
	if server {
		sPort := strconv.Itoa(port)
		if _, err := net.InterfaceByName(iface); err != nil {
			log.Fatalf("Could Not find interface %v: %v", iface, err)
		}
		cch := make(chan int32)
		go signalLoop(remote+":8000", cch)
		BandwidthTest(iface, sPort, fullAddr, snaplen, cch)
		LatencyTest(remote)
	} else {
		controlServer(remote+":8000", fullAddr)
	}
}
