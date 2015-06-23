package vdetesting

import (
	"net"
	"strconv"
)

//TestServer is the server side part of a test
//it should receive data and log all the statistic we need
type TestServer interface {
	StartServer()
	AddStat(s Stat)
	NetworkInterface() *net.Interface
}

//TestClient is the Client side part of a test
//it should send data and publish an method to use in a cycle
type TestClient interface {
	StartClient()
}

//Test is a generic test it need a client method and a server method
// and it test one single aspect and save the results to a single logfile
type Test interface {
	TestClient
	TestServer
	Name() string
	Address() net.Addr
	Port() Port
}

//Port is a Network Port that Contains the port number
//and the methods to use them
type Port struct {
	port int
}

func (p *Port) String() string {
	return strconv.Itoa(p.port)
}

//Int return the Integer for the Port
func (p *Port) Int() int {
	return p.port
}

//NextPort return you the next port in order
func (p *Port) NextPort(int) Port {
	next := Port{p.port + 1}
	return next
}

//Stat let you gather statistic regarding any kind of test
type Stat interface {
	Start()
	Stop()
}
