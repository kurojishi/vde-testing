package main

import (
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	kb int64 = 1000
	mb int64 = 1000 * kb
	gb int64 = 1000 * mb
)

var sizes = []int64{10, 30, 50, 100, 200, 500}

//controlServer start the controls channel on the client
func controlServer(bind, address string, port int) {
	clistener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Control Server Started")
	var stressCh []chan int32
	for {
		conn, err := clistener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		var buf int32
		binary.Read(conn, binary.LittleEndian, &buf)
		//TODO: define the other cases
		switch buf {
		case bandwidth:
			log.Print("Starting BandwidthTest")
			sendData(address+":"+strconv.Itoa(port), 1000)
		case die:
			break
		case stress:
			log.Print("Starting StressTest")
			stressCh = stressSend(address, port)
		case stressStop:
			log.Print("Stopping StressTest")
			for i := 0; i < len(stressCh); i++ {
				stressCh[i] <- stop
				close(stressCh[i])
			}
		//case load:
		default:
			continue
		}
	}
}

//sendData send size data (in megabytes)to the string addr
func sendData(addr string, size int64) {
	_, err := net.ResolveTCPAddr("tcp", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		//log.Printf("sendData: %v", err)
		return
	}
	//defer conn.Close()
	n, err := io.CopyN(conn, devZero, size*(mb))
	if err != nil {
		//log.Print(err)
		return
	}
	if n != size*mb {
		log.Printf("couldnt send %v Megabytes", float64(n)/float64(mb))
		return
	}
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
					sendData(finalAddr, sizes[r.Int31n(int32(len(sizes)-1))])
				}
			}
		}(finalAddr, ssch)

	}
	return controlChannels

}
