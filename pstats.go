package main

import (
	"log"
	"time"

	"github.com/jandre/procfs"
)

//Stats save stats gather information
//on on the pid process using the proc filesystem
func Stats(pid int, ticker *time.Ticker) error {
	for now := range ticker.C {
		process, err := procfs.NewProcess(pid, true)
		if err != nil {
			return err
		}
		pstats, err := process.Stat()
		if err != nil {
			return err
		}
		pstatsm, err := process.Statm()
		if err != nil {
			return err
		}
		pstatus, err := process.Status()
		if err != nil {
			return err
		}

		log.Printf("Polling vde_switch data %v: cputime: %v memory: %v context_switches: %v threads: %v", now, pstats.Utime.UnixNano()+pstats.Stime.UnixNano()/int64(time.Second), pstatsm.Size, pstatus.NVcswitch+pstatus.Vcswitch, pstats.NumThreads)
	}

	return nil
}

//PollStats call Stats every tick and close it when it's done
func PollStats(pid int) *time.Ticker {
	ticker := time.NewTicker(1 * time.Millisecond)
	go Stats(pid, ticker)
	return ticker
}
