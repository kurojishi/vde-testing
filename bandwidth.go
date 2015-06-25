// Package vdetesting provides the framework for testing
package vdetesting

import (
	"encoding/binary"
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
	cch     chan int32
	stats   StatManager
	kind    string
}

//NewBandwidthTest Return a new BandwidthTest
func NewBandwidthTest(kind string, iface string, address string, port int) (*BandwidthTest, error) {
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
	}
	bandwidth := BandwidthTest{iface: face,
		address: addr,
		port:    Port{port},
		cch:     make(chan int32),
		stats:   NewStatManager(),
		kind:    kind,
		name:    "bandwidth"}
	logfile := bandwidth.name + ".log"
	if kind == "server" {
		stat := NewTCPStat(bandwidth.iface, bandwidth.port, logfile)
		bandwidth.AddStat(&stat)
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

//StartServer start the server side of Bandwidthtest
func (t *BandwidthTest) StartServer() {
	log.Printf("Starting bandwidth test")
	err := t.stats.Start()
	if err != nil {
		log.Fatal(err)
	}
	bind, err := utils.InterfaceAddrv4(t.iface)
	listener, err := net.Listen("tcp", bind+":"+t.port.String())
	if err != nil {
		log.Fatalf("ReceiveData %v", err)
	}
	defer listener.Close()
	go utils.SendControlSignal(t.address.String(), 1)
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	utils.DevNullConnection(conn, nil)
	t.stats.Stop()
	log.Print("Finished bandwidth test")

}

//IFace Return the Interface
func (t *BandwidthTest) IFace() *net.Interface {
	return t.iface
}

//StartClient start the TestClient side of this Test
func (t *BandwidthTest) StartClient() error {
	var arrived = false
	local, err := utils.Localv4Addr()
	if err != nil {
		return err
	}
	clistener, err := net.Listen("tcp", local+":8999")
	log.Printf("waiting control message for bandwidth test on %v", local+":8999")
	if err != nil {
		log.Fatal(err)
	}
	for !arrived {
		conn, err := clistener.Accept()
		if err != nil {
			return err
		}
		var buf int32
		binary.Read(conn, binary.LittleEndian, &buf)
		if buf == 1 {
			arrived = true
			clistener.Close()
			log.Printf("control message arrived")
		} else {
			log.Print(buf)
		}
	}
	log.Print("control message arrived sending data")
	utils.SendData(t.address.String()+":"+t.port.String(), 1000)
	return nil
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
