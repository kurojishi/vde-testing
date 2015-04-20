package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendData(adr string) {
	var payload = []byte("skldjaslkcjlak")
	conn, err := net.Dial("tcp", adr)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("sending data")
	for i := 0; i < 20; i++ {
		_, err = conn.Write(payload)
	}
	conn.Close()
}

func readPackets(inter string, finished chan bool) {
	handle, err := pcap.OpenLive(inter, 1600, true, 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)
	for packet := range packetSource.Packets() {
		fmt.Println(packet.Dump())
	}
	finished <- true
}

func receiveData(conn net.Listener, ok chan bool) {
	fmt.Println("receiver starterd")
	payload := make([]byte, 14)
	listenerConnection, err := conn.Accept()
	if err != nil {
		fmt.Print(err)
	}
	for i := 0; i < 20; i++ {
		_, err := listenerConnection.Read(payload)
		fmt.Print(payload)
		if err != nil {
			fmt.Print(err)
		}
	}
	conn.Close()
	ok <- true

}

func main() {
	var in = flag.String("in", "192.168.4.1:5000", "address to send the data too")
	//var out = flag.String("out", "192.168.4.2:5000", "address to send the data from")
	//outAddress := net.ParseIP(*out)
	//if outAddress == nil {
	//log.Fatalln("bad output ip address")
	//}
	listener, err := net.Listen("tcp", *in)
	defer listener.Close()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("server started")
	ok := make(chan bool)
	finished := make(chan bool)
	go readPackets("tap0", finished)
	time.Sleep(time.Duration(1000))
	go receiveData(listener, ok)
	go sendData(*in)
	x := <-ok
	if !x {
		fmt.Println("WTF")
	}
	y := <-finished
	if !y {
		fmt.Println("WTF")
	}
}
