package main

import (
	"log"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

//TODO: save data to file
func latencyTest(address string) {
	rttch := make(chan time.Duration)
	for i := 0; i < 10; i++ {
		ra, err := net.ResolveIPAddr("ip", address)
		if err != nil {
			log.Fatal(err)
		}
		pinger := fastping.NewPinger()
		pinger.AddIPAddr(ra)
		pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
			rttch <- rtt

		}
		maxRtt, err := time.ParseDuration("1s")
		if err != nil {
			log.Fatal(err)
		}
		pinger.MaxRTT = maxRtt
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
