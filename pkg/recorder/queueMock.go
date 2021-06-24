package recorder

import "github.com/chifflier/nflog-go/nflog"

type QueueMock struct {
	initialized bool
	started     bool
	stopped     bool
}

func (q *QueueMock) Initialize(group int, handler nflog.Callback) error {
	q.initialized = true
	return nil
}

func (q *QueueMock) Start() error {
	q.started = true
	return nil
}

func (q *QueueMock) Stop() {
	q.stopped = true
}
