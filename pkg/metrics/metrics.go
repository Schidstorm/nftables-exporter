package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PacketCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "nftables_exporter_packet_count",
		Help: "The total number of network packets",
	}, []string {
		//name of the host
		"host",

		//group number from libnetfilter
		"group",

		//input interface name
		"iif",

		//output interface name
		"oif",

		//source ip
		"saddr",

		//destination port
		"dport",

		//is udp datagram
		"udp",

		//is tcp datagram
		"tcp",

		//protocol name
		"protocol",

		// ip version (4 or 6)
		"ipVersion",
	})
)

type PacketMetric struct {
	InputInterface  string
	OutputInterface string
	SourceIp        string
	DestinationIp   string
	DestinationPort uint16
	Udp             bool
	Tcp             bool
	Protocol        string
	IpVersion		int
}
