package vdetesting

import (
	"net"
	"strconv"
	"sync"
)

//TestServer is the server side part of a test
//it should receive data and log all the statistic we need
type TestServer interface {
	StartServer()
	AddStat(s Stat)
	statManager()
	IFace() *net.Interface
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
	SetWaitGroup(wg *sync.WaitGroup) error
}

//StatManager is a struct that should be added to everytest
//it manage all the Stats inside them
type StatManager struct {
	stats []Stat
	wg    *sync.WaitGroup
}

//NewStatManager Create a NewStatManager, should be used inside tests
func NewStatManager() StatManager {
	var wg sync.WaitGroup
	manager := StatManager{stats: make([]Stat, 0, 20),
		wg: &wg}
	return manager
}

//Add new statistic fetcher to the manager
func (manager *StatManager) Add(s Stat) {
	s.SetWaitGroup(manager.wg)
	manager.stats = append(manager.stats, s)
}

//Stats return the slice with all the Stats we fetch
func (manager *StatManager) Stats() *[]Stat {
	return &manager.stats
}

//Start start all the statistics
func (manager *StatManager) Start() error {
	for i := 0; i < len(manager.stats); i++ {
		manager.stats[i].Start()
	}
	return nil
}

//Stop stop all the statistics and wait for them to finish
//TODO: add waitgroup handling
func (manager *StatManager) Stop() error {
	for i := 0; i < len(manager.stats); i++ {
		manager.stats[i].Stop()
	}
	return nil
}
