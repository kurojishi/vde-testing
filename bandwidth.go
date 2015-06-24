// Package vdetesting provides the framework for testing
package vdetesting

import (
	"log"
	"net"
	"strconv"

	"github.com/kurojishi/vdetesting/utils"
)

type BandwidthTest struct {
	iface   *net.Interface
	address net.Addr
	port    Port
	name    string
	cch     chan int32
	stats   StatManager
}

func NewBandwidthTest(iface string, address string, port int) (*BandwidthTest, error) {
	addr, err := net.ResolveIPAddr("tcp", address+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	face, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	bandwidth := BandwidthTest{iface: face,
		address: addr,
		port:    Port{port},
		cch:     make(chan int32),
		stats:   StatManager{stats: make([]Stat, 0, 20)},
		name:    "bandwidth"}
	return &bandwidth, nil
}

func (t *BandwidthTest) AddStat(stat Stat) {
	t.stats.Add(stat)

}

func (t *BandwidthTest) statisticsStart() {
	t.stats.Start()

}

func (t *BandwidthTest) statisticsStop() {
	t.stats.Stop()
}

//StartServer start the server side of Bandwidthtest
func (t *BandwidthTest) StartServer() {
	log.Printf("Starting bandwidth test")
	//go TCPStats(iface, snaplen, port, sync)
	//ticker, sch := PollStats(pid, "bandwidth")
	listener, err := net.Listen("tcp", t.address.String())
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	defer listener.Close()
	t.cch <- bandwidth
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	utils.DevNullConnection(conn, nil)
	//ticker.Stop()
	//sch <- true
	//close(sch)
	log.Print("Finished bandwidth test")

}

func (t *BandwidthTest) IFace() *net.Interface {
	return t.iface
}

func (t *BandwidthTest) StartClient() {
	utils.SendData(t.address.String(), 1000)
}

func (t *BandwidthTest) Name() string {
	return t.name
}

func (t *BandwidthTest) Port() Port {
	return t.port
}

func (t *BandwidthTest) Address() net.Addr {
	return t.address
}
