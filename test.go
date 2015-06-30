package vdetesting

import (
	"errors"
	"log"
	"net"
	"sync"
)

//TestServer is the server side part of a test
//it should receive data and log all the statistic we need
type TestServer interface {
	StartServer()
	AddStat(s Stat)
	statManager()
	IFace() *net.Interface
	Name() string
	Address() net.Addr
}

//Test is a generic test it need a client method and a server method
// and it test one single aspect and save the results to a single logfile
type Test interface {
	StartClient()
	TestServer
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
	stats   []Stat
	wg      *sync.WaitGroup
	started bool
}

//NewStatManager Create a NewStatManager, should be used inside tests
func NewStatManager() StatManager {
	var wg sync.WaitGroup
	manager := StatManager{stats: make([]Stat, 0, 20),
		wg: &wg, started: false}
	return manager
}

//Add new statistic fetcher to the manager
func (manager *StatManager) Add(s Stat) {
	s.SetWaitGroup(manager.wg)
	manager.stats = append(manager.stats, s)
}

//Start start all the statistics
func (manager *StatManager) Start() error {
	if manager.started {
		err := errors.New("Statistics already Started")
		return err
	}
	for _, stat := range manager.stats {
		stat.Start()
	}
	manager.started = true
	return nil
}

//Stop stop all the statistics and wait for them to finish
func (manager *StatManager) Stop() error {
	for _, stat := range manager.stats {
		stat.Stop()
	}
	log.Print("waiting for stats to stop")
	manager.wg.Wait()
	manager.started = false
	return nil
}
