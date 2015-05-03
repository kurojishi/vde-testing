package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendData(addr string, ok chan bool) {
	var payload = []byte("skldjaslkcjlak")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("sending data")
	for i := 0; i < 20; i++ {
		_, err = conn.Write(payload)
	}
	conn.Close()
	ok <- true
}

func readPackets(inter string, ok chan bool) {
	handle, err := pcap.OpenLive(inter, 1600, true, 0)
	defer handle.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)
	packetSource.NoCopy = true
	fmt.Println("catcher ready")
	for {
		packet, err := packetSource.NextPacket()

		if err != nil {
			fmt.Print(err)
		}
		fmt.Println(string(packet.Data()))
	}
	ok <- true
}

func receiveData(conn net.Listener) {
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

}

func main() {
	var sender bool
	var in string
	flag.StringVar(&in, "in", "192.168.4.2:5000", "address to send the data too")
	flag.BoolVar(&sender, "server", false, "service will be a server")
	flag.Parse()
	if sender {
		listener, err := net.Listen("tcp", in)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println("server started")
		ok := make(chan bool)
		go readPackets("tap0", ok)
		if err != nil {
			fmt.Print(err)
		}
		go receiveData(listener)
		x := <-ok
		if !x {
			fmt.Println("WTF")
		}
	} else {
		ok := make(chan bool)
		go sendData(in, ok)
		x := <-ok
		if !x {
			fmt.Println("WTF")
		}
	}
}
