package main

import (
	"flag"
	"log"

	"github.com/kurojishi/vdetesting"
)

var pid int

func main() {
	var server bool
	var port int
	var remote, iface string
	flag.StringVar(&remote, "remote", "192.168.4.1", "address to send the data too")
	flag.IntVar(&port, "port", 5000, "starting port")
	flag.BoolVar(&server, "server", false, "service will be a server")
	flag.StringVar(&iface, "i", "tap0", "interface connected to the switch")
	flag.IntVar(&pid, "pid", 0, "the vde switch pid")
	flag.Parse()
	if server {
		//btest, err := vdetesting.NewBandwidthTest("server", iface, remote, port, pid)
		//if err != nil {
		//log.Fatal(err)
		//}
		//btest.StartServer()
		//logger, err := vdetesting.NewLatencyTest(iface, remote)
		//logger.StartServer()
		stest, err := vdetesting.NewStressTest("server", iface, remote, port, pid)
		if err != nil {
			log.Fatal(err)
		}
		stest.StartServer()
	} else {
		//btest, err := vdetesting.NewBandwidthTest("client", iface, remote, port, pid)
		//if err != nil {
		//log.Fatal(err)
		//}
		//btest.StartClient()
		stest, err := vdetesting.NewStressTest("client", iface, remote, port, pid)
		if err != nil {
			log.Fatal(err)
		}
		stest.StartClient()
	}

}
