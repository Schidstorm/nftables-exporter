package nftables_exporter

import (
	"github.com/chifflier/nflog-go/nflog"
	"github.com/sirupsen/logrus"
	"syscall"
)

type Queue struct {
	id    int
	queue *nflog.Queue
}

func NewQueue(group int, handler nflog.Callback) (*Queue, error) {

	q := new(nflog.Queue)
	queue := &Queue{
		id:    group,
		queue: q,
	}



	if err := q.SetCallback(handler); err != nil {
		return nil, err
	}

	if err := q.Init(); err != nil {
		return nil, err
	}

	if err := q.Unbind(syscall.AF_INET); err != nil {
		return nil, err
	}

	if err := q.Bind(syscall.AF_INET); err != nil {
		return nil, err
	}

	if err := q.CreateQueue(group); err != nil {
		return nil, err
	}

	if err := q.SetMode(nflog.NFULNL_COPY_PACKET); err != nil {
		return nil, err
	}

	logrus.Info("initialized queue")
	return queue, nil
}

// Start the queue.
func (q *Queue) Start() error {
	return q.queue.TryRun()
}

// Stop the queue.
func (q *Queue) Stop() {
	q.queue.Close()
}