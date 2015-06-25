package vdetesting

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/kurojishi/vdetesting/utils"
	"github.com/tatsushid/go-fastping"
)

const (
	bandwidth  int32 = 1
	latency    int32 = 2
	load       int32 = 3
	stress     int32 = 4
	stressStop int32 = 5
	die        int32 = 0
)

const (
	stop  int32 = 1
	ready int32 = 2
)

func signalLoop(control string, cch chan int32) {
	for msg := range cch {
		err := utils.SendControlSignal(control, msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//LatencyTest use ping to control latency
func LatencyTest(address string) {
	rttch := make(chan time.Duration, 10)
	ra, err := net.ResolveIPAddr("ip", address)
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
		log.Printf("Test ping:latency %v ms", (float32(rtt)/2)/float32(time.Millisecond))
	}
	log.Printf("Medium Latency is: %v", (float32(sum/2)/float32(i))/float32(time.Millisecond))

}

func manageConnections(address string, sch chan int32, wg sync.WaitGroup) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Manage Connections error: %v", err)
	}
	wg.Add(1)
	defer wg.Done()
	//defer listener.Close()
	for {
		select {
		case <-sch:
			return
		default:
			conn, err := listener.Accept()
			//defer conn.Close()
			if err != nil {
				log.Fatalf("Manage Connections error 2: %v", err)
			}
			go utils.DevNullConnection(conn, &wg)

		}
	}
	log.Printf("Closing connection %v", address)

}

//StressTest lauch a test to see what the vde_switch will do on very intensive traffic
func StressTest(address string, startingPort int, cch chan int32, pid int) {
	log.Print("Starting stress test")
	schContainer := make([]chan int32, 0, 50)
	//ticker, sch := PollStats(pid, "stress")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		ssch := make(chan int32, 1)
		finalAddr := address + ":" + strconv.Itoa(startingPort+i)
		go manageConnections(finalAddr, ssch, wg)
		schContainer = append(schContainer, ssch)

	}
	cch <- stress
	timer := time.NewTimer(1 * time.Minute)
	<-timer.C
	log.Print("timer elapsed")
	cch <- stressStop
	for i := 0; i < len(schContainer); i++ {
		schContainer[i] <- stop
	}
	//ticker.Stop()
	log.Print("Stopping ticker")
	//sch <- true
	log.Print("Asking Threads to stop")
	wg.Wait()
	log.Print("Finished stress test")
}
