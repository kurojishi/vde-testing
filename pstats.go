package vdetesting

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jandre/procfs"
)

//ProfilingStat implement Stat
//And it's used to fetch profiling data from /proc
type ProfilingStat struct {
	pid    int
	wg     *sync.WaitGroup
	sync   chan bool
	logger *log.Logger
	ticker *time.Ticker
}

//NewProfilingStat create  new stat for profiling a process using /proc
func NewProfilingStat(pid int, logfile string) ProfilingStat {
	logfile += "profiling.log"
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
	sync := make(chan bool, 1)
	ticker := time.NewTicker(500 * time.Millisecond)
	stat := ProfilingStat{pid: pid, logger: logger, sync: sync, ticker: ticker}
	return stat
}

//SetWaitGroup add a wait group to the Stat
//a wait group is used by the stat manager to correctly close and sync all the stats
//before returing control to the test
//return error if there is already a waitgroup setup
func (stat *ProfilingStat) SetWaitGroup(wg *sync.WaitGroup) error {
	if stat.wg != nil {
		return errors.New("there is already a wait group")
	}
	stat.wg = wg
	return nil
}

//Stats save stats gather information
//on on the pid process using the proc filesystem
func (stat *ProfilingStat) stats() {
	stat.wg.Add(1)
	for {
		select {
		case <-stat.ticker.C:
			process, err := procfs.NewProcess(stat.pid, true)
			if err != nil {
				log.Fatal(err)
			}
			pstats, err := process.Stat()
			if err != nil {
				log.Fatal(err)
			}
			pstatus, err := process.Status()
			if err != nil {
				log.Fatal(err)
			}
			stat.logger.Printf("%v %v %v %v", pstats.Utime+pstats.Stime, pstats.Vsize, pstatus.NVcswitch+pstatus.Vcswitch, pstats.NumThreads)
		case <-stat.sync:
			stat.wg.Done()
			return
		}
	}
	log.Println("Finished Polling Stats")
}

//Start start ProfilingStat
func (stat *ProfilingStat) Start() {
	go stat.stats()
}

//Stop send the signal to to the goroutine to stop
func (stat *ProfilingStat) Stop() {
	stat.ticker.Stop()
	stat.sync <- true
	close(stat.sync)
	log.Print("stopped profiling stat")
}
