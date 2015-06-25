package vdetesting

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/kurojishi/vdetesting/utils"
)

var sizes = []int64{10, 30, 50, 100, 200, 500}

func manageConnections(address string, sch chan bool, wg sync.WaitGroup) {
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
	schContainer := make([]chan bool, 0, 50)
	//ticker, sch := PollStats(pid, "stress")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		ssch := make(chan bool, 1)
		finalAddr := address + ":" + strconv.Itoa(startingPort+i)
		go manageConnections(finalAddr, ssch, wg)
		schContainer = append(schContainer, ssch)

	}
	timer := time.NewTimer(1 * time.Minute)
	<-timer.C
	log.Print("timer elapsed")
	for i := 0; i < len(schContainer); i++ {
		schContainer[i] <- true
	}
	//ticker.Stop()
	log.Print("Stopping ticker")
	//sch <- true
	log.Print("Asking Threads to stop")
	wg.Wait()
	log.Print("Finished stress test")
}

//composition of sendData
func stressSend(address string, startingPort int) []chan int32 {
	controlChannels := make([]chan int32, 0, 50)
	for i := 0; i < 50; i++ {
		ssch := make(chan int32)
		controlChannels = append(controlChannels, ssch)

		finalAddr := address + ":" + strconv.Itoa(startingPort+i)
		go func(finalAddr string, ssch chan int32) {

			for {
				select {
				case <-ssch:
					return
				default:
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					utils.SendData(finalAddr, sizes[r.Int31n(int32(len(sizes)-1))])
				}
			}
		}(finalAddr, ssch)

	}
	return controlChannels

}
