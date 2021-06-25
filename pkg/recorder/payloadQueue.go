package recorder

import (
	"sync"
)

type PayloadQueue struct {
	frontMutex sync.Mutex
	front uint32
	back uint32
	size uint32
	list []PayloadInterface
}

func NewPayloadQueue(size uint32) *PayloadQueue {
	return &PayloadQueue{
		front: 0,
		back:  0,
		size:  size,
		list:  make([]PayloadInterface, size),
	}
}

func (q *PayloadQueue) Enqueue(payload PayloadInterface)  {
	q.frontMutex.Lock()
	if q.back + 1 == q.front {
		q.frontMutex.Unlock()
		return
	}
	q.frontMutex.Unlock()


	q.list[q.back] = payload
	q.back = (q.back + 1) % q.size
}

func (q *PayloadQueue) Dequeue() (PayloadInterface, bool) {
	q.frontMutex.Lock()
	if q.back == q.front {
		q.frontMutex.Unlock()
		return nil, false
	}
	index := q.front
	q.front = (q.front + 1) % q.size
	q.frontMutex.Unlock()
	return q.list[index], true
}
