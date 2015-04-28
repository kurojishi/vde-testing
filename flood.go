package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendData(addr string, out string) {
	var payload = []byte("skldjaslkcjlak")
	localAddr, err := net.ResolveTCPAddr("tcp", out)
	if err != nil {
		fmt.Print(err)
	}
	dialer := net.Dialer{LocalAddr: localAddr}
	conn, err := dialer.Dial("tcp", addr)
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
	defer handle.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)
	packetSource.NoCopy = true
	fmt.Println("catcher ready")
	for i := 0; i < 20; i++ {
		packet, err := packetSource.NextPacket()

		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(packet.Dump())
		fmt.Println(i)
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
		fmt.Println(string(payload))
		if err != nil {
			fmt.Print(err)
		}
	}
	conn.Close()
	ok <- true

}

func main() {
	var in = flag.String("in", "192.168.4.1:5000", "address to send the data too")
	var out = flag.String("out", "192.168.4.15", "address to send the data from")
	var sender = flag.Bool("sender", false, "service will be a server")
	if !*sender {
		listener, err := net.Listen("tcp", *in)
		defer listener.Close()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println("server started")
		ok := make(chan bool)
		finished := make(chan bool)
		go readPackets("tap0", finished)
		if err != nil {
			fmt.Print(err)
		}
		go receiveData(listener, ok)
		y := <-finished
		if !y {
			fmt.Println("WTF")
		}
		x := <-ok
		if !x {
			fmt.Println("WTF")
		}
	} else {
		go sendData(*in, *out)
	}
}
