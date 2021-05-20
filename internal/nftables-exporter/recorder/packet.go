package recorder

import (
	"context"
	"github.com/chifflier/nflog-go/nflog"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schidstorm/nftables-exporter/internal/nftables-exporter/metrics"
	nftables_exporter "github.com/schidstorm/nftables-exporter/pkg/nftables-exporter"
	"strconv"
)

var handlerContext *handler

func Packet(ctx context.Context, group int, host string) error {
	handlerContext = &handler{
		group: group,
		host:  host,
	}

	if queue, err := nftables_exporter.NewQueue(group, Handle); err != nil {
		return err
	} else {
		for {
			if err := queue.Start(); err != nil {
				return err
			}

			select {
			case <-ctx.Done():
				queue.Stop()
				return ctx.Err()
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
	metric := ParsePacket(p)
	labels := prometheus.Labels{
		"udp":   "",
		"tcp":   "",
		"iif":   "",
		"oif":   "",
		"saddr": "",
		"dport": "",
		"group": strconv.Itoa(handlerContext.group),
		"host":  handlerContext.host,
		"protocol": "",
	}

	if metric.Udp {
		labels["udp"] = "1"
	} else {
		labels["udp"] = "0"
	}

	if metric.Tcp {
		labels["tcp"] = "1"
	} else {
		labels["tcp"] = "0"
	}

	if metric.InputInterface != nil {
		labels["iif"] = *metric.InputInterface
	}
	if metric.OutputInterface != nil {
		labels["oif"] = *metric.OutputInterface
	}
	if metric.SourceIp != nil {
		labels["saddr"] = *metric.SourceIp
	}

	if metric.DestinationPort != nil {
		labels["dport"] = strconv.Itoa(int(*metric.DestinationPort))
	}

	if metric.Protocol != nil {
		labels["protocol"] = *metric.Protocol
	}

	if labels["saddr"] != "0.0.0.0" {
		metrics.PacketCounter.With(labels).Inc()
	}

	return 0
}

func ParsePacket(payload *nflog.Payload) *metrics.PacketMetric {
	metric := &metrics.PacketMetric{}

	if inInterface, err := nftables_exporter.GetInterfaceFromNumber(payload.GetInDev()); err == nil {
		metric.InputInterface = new(string)
		*metric.InputInterface = inInterface.Attrs().Name
	}

	if outInterface, err := nftables_exporter.GetInterfaceFromNumber(payload.GetOutDev()); err == nil {
		metric.OutputInterface = new(string)
		*metric.OutputInterface = outInterface.Attrs().Name
	}

	packetV4 := gopacket.NewPacket(payload.Data, layers.LayerTypeIPv4, gopacket.Default)

	if ipLayer := packetV4.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ipv4, _ := ipLayer.(*layers.IPv4)
		metric.SourceIp = new(string)
		*metric.SourceIp = ipv4.SrcIP.String()
		metric.DestinationIp = new(string)
		*metric.DestinationIp = ipv4.DstIP.String()
		metric.Protocol = new(string)
		*metric.Protocol = ipv4.Protocol.String()
	}

	packetV6 := gopacket.NewPacket(payload.Data, layers.LayerTypeIPv6, gopacket.Default)
	if ipLayer := packetV6.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ipv6, _ := ipLayer.(*layers.IPv6)
		metric.SourceIp = new(string)
		*metric.SourceIp = ipv6.SrcIP.String()
		metric.DestinationIp = new(string)
		*metric.DestinationIp = ipv6.DstIP.String()
		metric.Protocol = new(string)
		*metric.Protocol = ipv6.NextHeader.String()
	}

	if tcpLayer := packetV4.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		metric.Tcp = true
		tcp, _ := tcpLayer.(*layers.TCP)
		metric.DestinationPort = new(uint16)
		*metric.DestinationPort = uint16(tcp.DstPort)
	} else {
		metric.Tcp = false
		if udpLayer := packetV4.Layer(layers.LayerTypeUDP); udpLayer != nil {
			metric.Udp = true
			udp, _ := udpLayer.(*layers.UDP)
			metric.DestinationPort = new(uint16)
			*metric.DestinationPort = uint16(udp.DstPort)
		} else {
			metric.Udp = false
		}
	}

	return metric
}
