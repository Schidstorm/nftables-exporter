package recorder

import (
	"context"
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schidstorm/nftables-exporter/pkg/metrics"
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
			"ipVersion": string(rune(metric.IpVersion)),
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