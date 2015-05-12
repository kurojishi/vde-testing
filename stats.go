package main

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

var statsResults []StatsStream

// simpleStreamFactory implements tcpassembly.StreamFactory
type statsStreamFactory struct{}

// StatsStream will handle the actual decoding of stats requests.
type StatsStream struct {
	net, transport                      gopacket.Flow
	bytes, packets, outOfOrder, skipped int64
	start, end                          time.Time
	sawStart, sawEnd                    bool
}

// New creates a new stream.  It's called whenever the assembler sees a stream
// it isn't currently following.
func (factory *statsStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	log.Printf("new stream %v:%v started", net, transport)
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
	statsResults = append(statsResults, *s)
}

//StreamStats returns all the statistics from a series of streams on a specific interface
// iface is the network interface to sniff and snaplen is the window size
func StreamStats(iface string, snaplen int32) []StatsStream {
	handle, err := pcap.OpenLive(iface, snaplen, true, 0)
	if err != nil {
		log.Fatal("error opening pcap handle: ", err)
	}
	defer handle.Close()

	//set up assembler

	streamFactory := &statsStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("Catching stream stats")

	var eth layers.Ethernet
	var dot1q layers.Dot1Q
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var ip6ext layers.IPv6ExtensionSkipper
	var tcp layers.TCP
	var udp layers.UDP
	var payload gopacket.Payload
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &dot1q, &ip4, &ip6, &ip6ext, &tcp, &udp, &payload)
	decoded := make([]gopacket.LayerType, 0, 4)

	var byteCount int64

	for {
		data, ci, err := handle.ReadPacketData()
		if err != nil {
			log.Printf("error getting packet: %v", err)
			//TODO:comunicate with client and restart
			continue
		}
		err = parser.DecodeLayers(data, &decoded)
		if err != nil {
			log.Printf("error decoding packet: %v", err)
			continue
		}

		byteCount += int64(len(data))
		foundNetLayer := false
		var netFlow gopacket.Flow
		for _, typ := range decoded {
			switch typ {
			case layers.LayerTypeIPv4:
				netFlow = ip4.NetworkFlow()
				foundNetLayer = true
			case layers.LayerTypeIPv6:
				netFlow = ip6.NetworkFlow()
				foundNetLayer = true
			case layers.LayerTypeTCP:
				if foundNetLayer {
					assembler.AssembleWithTimestamp(netFlow, &tcp, ci.Timestamp)
				}
			}
		}
	}
	assembler.FlushAll()
	return statsResults
}
