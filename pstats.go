package main

import (
	"log"
	"os"
	"time"

	"github.com/jandre/procfs"
)

//Stats save stats gather information
//on on the pid process using the proc filesystem
func Stats(pid int, ticker *time.Ticker, stop chan bool, testname string) error {
	if _, err := os.Stat(testname + "stats"); err == nil {
		err := os.Remove(testname + "stats")
		if err != nil {
			log.Fatal(err)
		}
	}
	logfile, err := os.Create(testname + "stats")
	defer log.Println("Finished Polling Stats")
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logfile, "", 0)
	for {
		select {
		case <-ticker.C:

			process, err := procfs.NewProcess(pid, true)
			if err != nil {
				return err
			}
			pstats, err := process.Stat()
			if err != nil {
				return err
			}
			pstatus, err := process.Status()
			if err != nil {
				return err
			}
			logger.Printf("%v %v %v %v", pstats.Utime+pstats.Stime, pstats.Vsize, pstatus.NVcswitch+pstatus.Vcswitch, pstats.NumThreads)
		case <-stop:
			return nil

		}
	}

	return nil
}

//PollStats call Stats every tick and close it when it's done
func PollStats(pid int, testname string) (*time.Ticker, chan bool) {
	stop := make(chan bool, 1)
	ticker := time.NewTicker(500 * time.Millisecond)
	go Stats(pid, ticker, stop, testname)
	return ticker, stop
}
