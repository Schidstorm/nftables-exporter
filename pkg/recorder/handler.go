package recorder

import (
	"context"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schidstorm/nftables-exporter/pkg/metrics"
	nftables_exporter "github.com/schidstorm/nftables-exporter/pkg/nftables-exporter"
	"github.com/sirupsen/logrus"
	"strconv"
)

type Handler struct {

}

func CreateHandlers(ctx context.Context, payloadQueue goconcurrentqueue.Queue, count int) {
	for i := 0; i < count; i++ {
		go CreateHandler(ctx, payloadQueue)
	}
}

func CreateHandler(ctx context.Context, payloadQueue goconcurrentqueue.Queue) {
	for {
		payload, err := payloadQueue.DequeueOrWaitForNextElementContext(ctx)
		if err != nil {
			logrus.Error(err)
			continue
		}

		metric := ParsePayload(payload.(PayloadInterface))
		if metric == nil {
			return
		}

		labels := prometheus.Labels{
			"udp":      "",
			"tcp":      "",
			"iif":      metric.InputInterface,
			"oif":      metric.OutputInterface,
			"saddr":    metric.SourceIp,
			"dport":    "",
			"ipVersion": strconv.Itoa(metric.IpVersion),
			"group":    strconv.Itoa(handlerContext.group),
			"host":     handlerContext.host,
			"protocol": metric.Protocol,
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

		if metric.DestinationPort != 0 {
			labels["dport"] = strconv.Itoa(int(metric.DestinationPort))
		}

		metrics.PacketCounter.With(labels).Inc()


	}
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
