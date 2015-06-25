package vdetesting

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

// simpleStreamFactory implements tcpassembly.StreamFactory
type statsStreamFactory struct {
	logger *log.Logger
}

// StatsStream will handle the actual decoding of stats requests.
type StatsStream struct {
	net, transport                      gopacket.Flow
	bytes, packets, outOfOrder, skipped int64
	start, end                          time.Time
	sawStart, sawEnd                    bool
	logger                              *log.Logger
}

var finished bool

// New creates a new stream.  It's called whenever the assembler sees a stream
// it isn't currently following.
func (factory *statsStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	//TODO:remove this print
	s := &StatsStream{
		net:       net,
		transport: transport,
		start:     time.Now(),
		logger:    factory.logger,
	}
	s.end = s.start
	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return s
}

// Reassembled is called whenever new packet data is available for reading.
// Reassembly objects contain stream data IN ORDER.
func (s *StatsStream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		if reassembly.Seen.Before(s.end) {
			s.outOfOrder++
		} else {
			s.end = reassembly.Seen
		}
		s.bytes += int64(len(reassembly.Bytes))
		s.packets++
		if reassembly.Skip > 0 {
			s.skipped += int64(reassembly.Skip)
		}
		s.sawStart = s.sawStart || reassembly.Start
		s.sawEnd = s.sawEnd || reassembly.End
	}
}

// ReassemblyComplete is called when the TCP assembler believes a stream has
// finished.
func (s *StatsStream) ReassemblyComplete() {
	if !finished {
		diffSecs := float64(s.end.Sub(s.start)) / float64(time.Second)
		s.logger.Printf("%v %v %v", diffSecs, float64(s.bytes)/float64(1000000), (float64(s.bytes)/float64(1000000))/diffSecs)
		log.Printf("%v %v %v", diffSecs, float64(s.bytes)/float64(1000000), (float64(s.bytes)/float64(1000000))/diffSecs)
	}
}

//TCPStat is a stat implementation
//for getting tcp statistic
type TCPStat struct {
	iface   *net.Interface
	port    Port
	sync    chan bool
	snaplen int
	logger  *log.Logger
	wg      *sync.WaitGroup
}

//NewTCPStat create  new tcp stat
func NewTCPStat(iface *net.Interface, port Port, logfile string) TCPStat {
	if _, err := os.Stat(logfile); err == nil {
		err := os.Remove(logfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Create(logfile)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", 0)
	sync := make(chan bool, 1)
	stat := TCPStat{iface: iface, port: port, logger: logger, snaplen: 1600, sync: sync}
	return stat
}

//SetWaitGroup add a wait group to the Stat
//a wait group is used by the stat manager to correctly close and sync all the stats
//before returing control to the test
//return error if there is already a waitgroup setup
func (s *TCPStat) SetWaitGroup(wg *sync.WaitGroup) error {
	if s.wg != nil {
		return errors.New("there is already a wait group")
	}
	s.wg = wg
	return nil
}

//Stop send the signal to the stat manager to stop polling stats
func (s *TCPStat) Stop() {
	s.sync <- true
	close(s.sync)
}

//Start returns all the statistics from a series of streams on a specific interface
// iface is the network interface to sniff and snaplen is the window size
func (s *TCPStat) Start() {
	go s.ifacePoll()
}

func (s TCPStat) ifacePoll() {
	finished = false
	flushDuration, err := time.ParseDuration("1m")
	if err != nil {
		log.Fatal("invalid flush duration", err)
	}
	log.Printf("starting capture on %v", s.iface.Name)

	//set up assembler

	streamFactory := &statsStreamFactory{s.logger}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)
	assembler.MaxBufferedPagesTotal = 0
	assembler.MaxBufferedPagesPerConnection = 0
	defer assembler.FlushAll()

	var eth layers.Ethernet
	var dot1q layers.Dot1Q
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var ip6ext layers.IPv6ExtensionSkipper
	var tcp layers.TCP
	var payload gopacket.Payload
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &dot1q, &ip4, &ip6, &ip6ext, &tcp, &payload)
	decoded := make([]gopacket.LayerType, 0, 8)

	var byteCount int64

	handle, err := pcap.OpenLive(s.iface.Name, int32(s.snaplen), true, flushDuration/2)
	if err != nil {
		log.Fatal("error opening pcap handle: ", err)
	}
	defer handle.Close()
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true
	nextFlush := time.Now().Add(flushDuration / 2)

	log.Println("Catching stream stats")
	s.wg.Add(1)
	for !finished {
		if time.Now().After(nextFlush) {
			assembler.FlushOlderThan(time.Now().Add(flushDuration))
			nextFlush = time.Now().Add(flushDuration / 2)
		}
		packet, err := source.NextPacket()
		if err != nil {
			continue
		}
		if err := parser.DecodeLayers(packet.Data(), &decoded); err != nil {
			log.Printf("error decoding packet: %v", err)
			continue
		}

		byteCount += int64(len(packet.Data()))
		if packet.TransportLayer().TransportFlow().Dst().String() == s.port.String() {
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), &tcp, packet.Metadata().Timestamp)
		}
	}
	log.Println("OUT OF THERE")
	<-s.sync
	s.wg.Done()
	log.Print("Catching finished")
}
