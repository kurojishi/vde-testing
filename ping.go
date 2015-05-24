package main

import (
	"log"
	"net"
	"time"

	"github.com/tatsushid/go-fastping"
)

func latencyTest(address string) {
	rttch := make(chan time.Duration, 10)
	rttContainer := make([]time.Duration, 0, 10)
	pinger := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ipv4:icmp", address)
	if err != nil {
		log.Fatal(err)
	}
	pinger.AddIPAddr(ra)
	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		rttch <- rtt
	}
	pinger.MaxRTT, _ = time.ParseDuration("3s")
	pinger.RunLoop()

	select {
	case <-pinger.Done():
		if err := pinger.Err(); err != nil {
			log.Fatal(err)
		}
	case rtt := <-rttch:
		rttContainer = append(rttContainer, rtt)
		if len(rttContainer) >= 10 {
			pinger.Stop()
		}
	}

	var sum time.Duration
	for value := range rttContainer {
		log.Printf("Test pint %d, latency %d", value, rttContainer[value]/2)
		sum += rttContainer[value] / 2
	}
	log.Printf("Medium latency is %v", sum/time.Second)

}
