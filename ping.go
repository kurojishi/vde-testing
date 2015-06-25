package vdetesting

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

//LatencyTest is a server side only
//Test that use ping to test latency
type LatencyTest struct {
	address net.Addr
	name    string
	iface   *net.Interface
	logger  *log.Logger
}

//NewLatencyTest Return a new LatencyTest
func NewLatencyTest(iface string, address string) (*LatencyTest, error) {
	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return nil, err
	}
	var face *net.Interface
	face, err = net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	logfile := "latency.log"
	if _, err := os.Stat(logfile); err == nil {
		err := os.Remove(logfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create(logfile)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", 0)
	latency := LatencyTest{iface: face,
		address: addr,
		name:    "latency",
		logger:  logger}
	return &latency, nil
}

//IFace Return the Interface
func (t *LatencyTest) IFace() *net.Interface {
	return t.iface
}

//Name return the name of this test
func (t *LatencyTest) Name() string {
	return t.name
}

//Address return the IP address of the test
func (t *LatencyTest) Address() net.Addr {
	return t.address
}

//StartServer use ping to control latency
func (t *LatencyTest) StartServer() {
	rttch := make(chan time.Duration, 10)
	ra, err := net.ResolveIPAddr("ip", t.address.String())
	if err != nil {
		log.Fatal(err)
	}
	pinger := fastping.NewPinger()
	pinger.AddIPAddr(ra)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rttch <- rtt

	}
	//set message size to 64 byte
	pinger.Size = 64
	for i := 0; i < 10; i++ {
		pinger.Run()
	}
	close(rttch)

	var sum time.Duration
	var i int
	for rtt := range rttch {
		sum += rtt
		i++
		t.logger.Printf("Test ping:latency %v ms", (float32(rtt)/2)/float32(time.Millisecond))
	}
	t.logger.Printf("Medium Latency is: %v", (float32(sum/2)/float32(i))/float32(time.Millisecond))

}
