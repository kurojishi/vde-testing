package vdetesting

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

// simpleStreamFactory implements tcpassembly.StreamFactory
type statsStreamFactory struct{}

// StatsStream will handle the actual decoding of stats requests.
type StatsStream struct {
	net, transport                      gopacket.Flow
	bytes, packets, outOfOrder, skipped int64
	start, end                          time.Time
	sawStart, sawEnd                    bool
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
	//TODO: write this data to a file
	if !finished {
		diffSecs := float64(s.end.Sub(s.start)) / float64(time.Second)
		log.Printf("stream took %v seconds to complete, sent %v MB with a bitrate of %v MB", diffSecs, float64(s.bytes)/float64(1000000), (float64(s.bytes)/float64(1000000))/diffSecs)
		finished = true
	}
}

//TCPStat is a stat implementation
//for getting tcp statistic
type TCPStat struct {
	iface   string
	port    Port
	sync    chan bool
	snaplen int
	logfile string
	wg      *sync.WaitGroup
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
	log.Printf("starting capture on %v", s.iface)

	//set up assembler

	streamFactory := &statsStreamFactory{}
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
	//var udp layers.UDP
	var payload gopacket.Payload
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &dot1q, &ip4, &ip6, &ip6ext, &tcp, &payload)
	decoded := make([]gopacket.LayerType, 0, 8)

	var byteCount int64

	handle, err := pcap.OpenLive(s.iface, int32(s.snaplen), true, flushDuration/2)
	if err != nil {
		log.Fatal("error opening pcap handle: ", err)
	}
	defer handle.Close()
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true
	nextFlush := time.Now().Add(flushDuration / 2)

	log.Println("Catching stream stats")
	s.sync <- true
	for !finished {
		if time.Now().After(nextFlush) {
			//log.Println("Flushing all streams that havent' seen packets")
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
	s.sync <- true
	log.Print("Catching finished")
}
