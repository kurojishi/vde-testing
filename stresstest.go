package vdetesting

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/kurojishi/vdetesting/utils"
)

var sizes = []int64{10, 30, 50, 100, 200, 500}

//StressTest is used for profiling and it's main use
//is to check how vde works under heavy loads
type StressTest struct {
	name    string
	iface   *net.Interface
	stats   StatManager
	kind    string
	pid     int
	port    Port
	address net.Addr
	sync    chan bool
}

//NewStressTest Return a new StressTest
func NewStressTest(kind string, iface string, address string, port int, pid int) (*StressTest, error) {
	addr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return nil, err
	}
	var face *net.Interface
	if kind == "server" {
		face, err = net.InterfaceByName(iface)
		if err != nil {
			return nil, err
		}
	} else {
		face = nil
		pid = 0
	}
	stress := StressTest{iface: face,
		address: addr,
		port:    Port{port},
		stats:   NewStatManager(),
		kind:    kind,
		pid:     pid,
		name:    "stress"}
	if kind == "server" {
		//tcpStat := NewTCPStat(stress.iface, stress.port, stress.name)
		//stress.AddStat(&tcpStat)
		pstat := NewProfilingStat(stress.pid, stress.name)
		stress.AddStat(&pstat)
	}
	return &stress, nil
}

//AddStat Add a new Statistic
func (t *StressTest) AddStat(stat Stat) {
	t.stats.Add(stat)

}

func (t *StressTest) statisticsStart() {
	t.stats.Start()

}

func (t *StressTest) statisticsStop() {
	t.stats.Stop()
}

//IFace Return the Interface
func (t *StressTest) IFace() *net.Interface {
	return t.iface
}

//Name return the name of this test
func (t *StressTest) Name() string {
	return t.name
}

//Port return the port this test will be performed on
func (t *StressTest) Port() Port {
	return t.port
}

//Address return the IP address of the test
func (t *StressTest) Address() net.Addr {
	return t.address
}

//manageConnection open all the connections on a single port
func (t *StressTest) manageConnections(address string, port Port, sch chan bool, listenGroup sync.WaitGroup) {
	listener, err := net.Listen("tcp", address+":"+port.String())
	if err != nil {
		log.Fatalf("Manage Connections error: %v", err)
	}
	listenGroup.Add(1)
	defer listenGroup.Done()
	for {
		select {
		case <-sch:
			log.Printf("i'm in there motherfuckers")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Fatalf("Manage Connections error 2: %v", err)
			}
			err = utils.DevNullConnection(conn, &listenGroup)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
	log.Printf("Closing connection %v", address)

}

//StartServer lauch a test to see what the vde_switch will do on very intensive traffic
func (t *StressTest) StartServer() {
	log.Print("Starting stress test")
	schContainer := make([]chan bool, 0, 50)
	var listenGroup sync.WaitGroup
	t.statisticsStart()
	laddr, err := utils.InterfaceAddrv4(t.iface)
	for i := 0; i < 50; i++ {
		ssch := make(chan bool, 1)
		go t.manageConnections(laddr, t.port.NextPort(i), ssch, listenGroup)
		schContainer = append(schContainer, ssch)

	}
	log.Println("Send start signal to client")
	err = utils.SendControlSignal(t.address.String(), 2)
	if err != nil {
		log.Fatal(err)
	}
	timer := time.NewTimer(1 * time.Minute)
	<-timer.C
	err = utils.SendControlSignal(t.address.String(), 2)
	log.Println("sent stop message to client waiting for confirmation")
	err = utils.WaitForControlMessage(laddr, 2)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(schContainer); i++ {
		schContainer[i] <- true
	}
	listenGroup.Wait()
	t.statisticsStop()
	log.Print("Finished stress test")
}
func (t *StressTest) sendRandomAmoutOfData(finalAddr string, ssch chan bool, sendGroup sync.WaitGroup) {
	sendGroup.Add(1)
	defer sendGroup.Done()
	for {
		select {
		case <-ssch:
			log.Printf("I'm here stopping %v goroutine", finalAddr)
			return
		default:
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			err := utils.SendData(finalAddr, sizes[r.Int31n(int32(len(sizes)-1))])
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

//StartClient is a composition of sendData
func (t *StressTest) StartClient() {
	var sendGroup sync.WaitGroup
	local, err := utils.Localv4Addr()
	if err != nil {
		log.Fatal(err)
	}
	err = utils.WaitForControlMessage(local, 2)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("control message arrived sending data")
	controlChannels := make([]chan bool, 0, 50)
	for i := 0; i < 50; i++ {
		ssch := make(chan bool, 1)
		controlChannels = append(controlChannels, ssch)
		port := t.port.NextPort(i)
		finalAddr := t.address.String() + ":" + port.String()
		go t.sendRandomAmoutOfData(finalAddr, ssch, sendGroup)
	}
	log.Println("Waiting for stop message")
	err = utils.WaitForControlMessage(local, 2)
	if err != nil {
		log.Fatal(err)
	}
	for _, cch := range controlChannels {
		cch <- true
		close(cch)
	}
	sendGroup.Wait()
	log.Println("Sending stop signal")
	err = utils.SendControlSignal(t.address.String(), 2)
	if err != nil {
		log.Fatal(err)
	}
}
