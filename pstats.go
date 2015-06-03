package main

import (
	"log"

	"github.com/jandre/procfs"
)

//PollVdeStats save stats gather information
//on on the pid process using the proc filesystem
func PollVdeStats(pid int, stop chan bool) error {
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

	log.Printf("Polling vde_switch data: cputime: %v memory: %v context_switches: %v threads: %v", pstats.Utime.UnixNano()+pstats.Stime.UnixNano(), pstatsm.Size, pstatus.NVcswitch+pstatus.Vcswitch, pstats.NumThreads)

	return nil
}
