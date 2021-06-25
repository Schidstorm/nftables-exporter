package recorder

import (
	"context"
	"github.com/chifflier/nflog-go/nflog"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/schidstorm/nftables-exporter/pkg/metrics"
	nftables_exporter "github.com/schidstorm/nftables-exporter/pkg/nftables-exporter"
	"github.com/sirupsen/logrus"
)

var handlerContext *handler
var payloadQueue = goconcurrentqueue.NewFIFO()
var MaxEnqueueRetriesCount = 1000

func init() {
	CreateHandlers(context.Background(), payloadQueue, 4)
}

func Packet(ctx context.Context, group int, host string, queueInterface nftables_exporter.QueueInterface) error {
	handlerContext = &handler{
		group: group,
		host:  host,
	}

	queue := queueInterface

	if err := queue.Initialize(group, Handle); err != nil {
		return err
	} else {
		for {
			if err := queue.Start(); err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				queue.Stop()
				return nil
			default:
				continue
			}
		}
	}
}

type handler struct {
	group int
	host  string
}

func Handle(p *nflog.Payload) int {
	HandlePayload(&Payload{nflogPayload: p})
	return 0
}

func HandlePayload(p PayloadInterface) int {
	for i := 0; i < MaxEnqueueRetriesCount; i++ {
		if err := payloadQueue.Enqueue(p); err == nil {
			break
		}
	}

	return 0
}

func ParsePayload(payload PayloadInterface) *metrics.PacketMetric {
	metric := &metrics.PacketMetric{}

	if inInterface, err := nftables_exporter.GetInterfaceFromNumber(payload.GetInDev()); err == nil {
		metric.InputInterface = inInterface.Attrs().Name
	} else {
		logrus.Error(err)
	}

	if outInterface, err := nftables_exporter.GetInterfaceFromNumber(payload.GetOutDev()); err == nil {
		metric.OutputInterface = outInterface.Attrs().Name
	} else {
		logrus.Error(err)
	}

	packetV4 := gopacket.NewPacket(payload.GetData(), layers.LayerTypeIPv4, gopacket.NoCopy)
	if ipLayer := packetV4.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ipv4, _ := ipLayer.(*layers.IPv4)
		if ipv4.Version == 4 {
			metric.IpVersion = 4
			metric.SourceIp = ipv4.SrcIP.String()
			metric.DestinationIp = ipv4.DstIP.String()
			metric.Protocol = ipv4.Protocol.String()
		}
	}

	packetV6 := gopacket.NewPacket(payload.GetData(), layers.LayerTypeIPv6, gopacket.NoCopy)
	if ipLayer := packetV6.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ipv6, _ := ipLayer.(*layers.IPv6)
		if ipv6.Version == 6 {
			metric.IpVersion = 4
			metric.SourceIp = ipv6.SrcIP.String()
			metric.DestinationIp = ipv6.DstIP.String()
			metric.Protocol = ipv6.NextHeader.String()
		}
	}

	if metric.IpVersion == 0 {
		return nil
	}

	if tcpLayer := packetV4.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		metric.Tcp = true
		tcp, _ := tcpLayer.(*layers.TCP)
		metric.DestinationPort = uint16(tcp.DstPort)
	} else {
		metric.Tcp = false
		if udpLayer := packetV4.Layer(layers.LayerTypeUDP); udpLayer != nil {
			metric.Udp = true
			udp, _ := udpLayer.(*layers.UDP)
			metric.DestinationPort = uint16(udp.DstPort)
		} else {
			metric.Udp = false
		}
	}

	return metric
}
