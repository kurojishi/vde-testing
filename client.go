package vdetesting

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/kurojishi/vdetesting/utils"
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
			utils.SendData(address+":"+strconv.Itoa(port), 1000)
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
