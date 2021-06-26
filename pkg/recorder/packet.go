package recorder

import (
	"context"
	"github.com/chifflier/nflog-go/nflog"
	"github.com/enriquebris/goconcurrentqueue"
	nftables_exporter "github.com/schidstorm/nftables-exporter/pkg/nftables-exporter"
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
	HandlePayload(&Payload{
		nflogPayload: p,
		inDev: p.GetInDev(),
		outDev: p.GetOutDev(),
	})
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

