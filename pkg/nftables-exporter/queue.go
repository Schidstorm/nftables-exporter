package nftables_exporter

import (
	"errors"
	"github.com/chifflier/nflog-go/nflog"
	"github.com/sirupsen/logrus"
	"syscall"
)

type Queue struct {
	id          int
	queue       *nflog.Queue
	initialized bool
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Initialize(group int, handler nflog.Callback) error {
	if q.initialized {
		return errors.New("called initialize twice")
	}

	q.queue = new(nflog.Queue)
	q.id = group
	q.initialized = true

	if err := q.queue.SetCallback(handler); err != nil {
		return err
	}

	if err := q.queue.Init(); err != nil {
		return err
	}

	if err := q.queue.Unbind(syscall.AF_INET); err != nil {
		return err
	}

	if err := q.queue.Bind(syscall.AF_INET); err != nil {
		return err
	}

	if err := q.queue.CreateQueue(group); err != nil {
		return err
	}

	if err := q.queue.SetMode(nflog.NFULNL_COPY_PACKET); err != nil {
		return err
	}

	logrus.Info("initialized queue")

	return nil
}

// Start the queue.
func (q *Queue) Start() error {
	if !q.initialized {
		return errors.New("tried to start uninitialized queue")
	}

	return q.queue.TryRun()
}

// Stop the queue.
func (q *Queue) Stop() {
	q.queue.Close()
}
