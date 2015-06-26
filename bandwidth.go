// Package vdetesting provides the framework for testing
package vdetesting

import (
	"log"
	"net"

	"github.com/kurojishi/vdetesting/utils"
)

//BandwidthTest is a Test that check the bandwidth
//of  a connection
type BandwidthTest struct {
	iface   *net.Interface
	address net.Addr
	port    Port
	name    string
	stats   StatManager
	kind    string
	pid     int
}

//NewBandwidthTest Return a new BandwidthTest
func NewBandwidthTest(kind string, iface string, address string, port int, pid int) (*BandwidthTest, error) {
	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return nil, err
	}
	var face *net.Interface
	if kind == "server" {
		face, err = net.InterfaceByName(iface)
		if err != nil {
			return nil, err
		}
	} else {
		face = nil
		pid = 0
	}
	bandwidth := BandwidthTest{iface: face,
		address: addr,
		port:    Port{port},
		stats:   NewStatManager(),
		kind:    kind,
		pid:     pid,
		name:    "bandwidth"}
	if kind == "server" {
		tcpStat := NewTCPStat(bandwidth.iface, bandwidth.port, bandwidth.name)
		bandwidth.AddStat(&tcpStat)
		pstat := NewProfilingStat(bandwidth.pid, bandwidth.name)
		bandwidth.AddStat(&pstat)
	}
	return &bandwidth, nil
}

//AddStat Add a new Statistic
func (t *BandwidthTest) AddStat(stat Stat) {
	t.stats.Add(stat)

}

func (t *BandwidthTest) statisticsStart() {
	t.stats.Start()

}

func (t *BandwidthTest) statisticsStop() {
	t.stats.Stop()
}

//IFace Return the Interface
func (t *BandwidthTest) IFace() *net.Interface {
	return t.iface
}

//Name return the name of this test
func (t *BandwidthTest) Name() string {
	return t.name
}

//Port return the port this test will be performed on
func (t *BandwidthTest) Port() Port {
	return t.port
}

//Address return the IP address of the test
func (t *BandwidthTest) Address() net.Addr {
	return t.address
}

//StartServer start the server side of Bandwidthtest
func (t *BandwidthTest) StartServer() {
	log.Printf("Starting bandwidth test")
	bind, err := utils.InterfaceAddrv4(t.iface)
	listener, err := net.Listen("tcp", bind+":"+t.port.String())
	if err != nil {
		log.Fatalf("Could not bind to listener %v", err)
	}
	defer listener.Close()
	go utils.SendControlSignal(t.address.String(), 1)
	t.statisticsStart()
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	utils.DevNullConnection(conn, nil)
	t.statisticsStop()
	log.Print("Finished bandwidth test")

}

//StartClient start the TestClient side of this Test
func (t *BandwidthTest) StartClient() {
	err := utils.WaitForControlMessage(1)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("control message arrived sending data")
	utils.SendData(t.address.String()+":"+t.port.String(), 1000)
}
