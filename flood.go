package main

import (
	"flag"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendData(addr string, ok chan bool) {
	payload := make([]byte, 4096)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sending data")
	for i := 0; i < 20; i++ {
		_, err = conn.Write(payload)
	}
	conn.Close()
	ok <- true
}

func packetSniffer(inter string, ok chan bool) {
	handle, err := pcap.OpenLive(inter, 65536, true, 0)
	defer handle.Close()
	if err != nil {
		log.Println(err.Error())
	}
	packetSource := gopacket.NewPacketSource(handle, layers.LinkTypeEthernet)
	packetSource.NoCopy = true
	log.Println("catcher ready")
	for packet := range packetSource.Packets() {
		if err != nil {
			log.Fatal(err)
		}
		log.Println(packet.Dump())
	}
	ok <- true
}

func receiveData(conn net.Listener) {
	log.Println("receiver starterd")
	payload := make([]byte, 4096)
	listenerConnection, err := conn.Accept()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		_, err := listenerConnection.Read(payload)
		log.Println(string(payload))
		if err != nil {
			log.Fatal(err)
		}
	}
	conn.Close()

}

func main() {
	var sender bool
	var address string
	var port string
	flag.StringVar(&address, "address", "192.168.4.1", "address to send the data too")
	flag.StringVar(&port, "p", "5000", "starting port")
	flag.BoolVar(&sender, "server", false, "service will be a server")
	flag.Parse()
	if sender {
		listener, err := net.Listen("tcp", address+":"+port)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("server started")
		ok := make(chan bool)
		go packetSniffer("tap0", ok)
		if err != nil {
			log.Fatal(err)
		}
		go receiveData(listener)
		x := <-ok
		if !x {
			log.Fatal(err)
		}
	} else {
		ok := make(chan bool)
		go sendData(address+":"+port, ok)
		x := <-ok
		if !x {
			log.Println("WTF")
		}
	}
}
