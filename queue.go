package queue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var emptyQueue = func() error {
	return errors.New("empty queue")
}

var fullQueue = errors.New("full queue")

type Queue struct {
	sync.RWMutex
	fs              []func() error
	stop            chan struct{}
	max, head, tail int
	tick            time.Duration
}

func NewQueue(capacity int, tick time.Duration) *Queue {
	return &Queue{
		fs:   make([]func() error, capacity),
		stop: make(chan struct{}),
		max:  capacity + 1,
		tick: tick,
	}
}

func (dq *Queue) Run() {
	ticker := time.NewTicker(dq.tick)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f := dq.Pop()
			go f.Do()
		case <-dq.stop:
			close(dq.stop)
			dq.Lock()
			dq.fs, dq.head, dq.tail = nil, 0, 0
			dq.Unlock()
			return
		}
	}
}

func (dq *Queue) Close() {
	defer func() {
		dq.stop <- struct{}{}
	}()
}

func (dq *Queue) Put(f Executor) error {
	dq.Lock()
	defer dq.Unlock()
	if dq.isFull() {
		return fullQueue
	}
	dq.fs[dq.tail] = f
	dq.tail = (dq.tail + 1) % dq.max
	return nil
}

func (dq *Queue) Pop() Executor {
	dq.Lock()
	defer dq.Unlock()
	if dq.isEmpty() {
		return emptyQueue
	}
	f := dq.fs[dq.head]
	dq.head = (dq.head + 1) % dq.max
	return f
}

func (dq *Queue) Size() int {
	dq.RLock()
	defer dq.RUnlock()
	return dq.size()
}

func (dq *Queue) size() int {
	return (dq.tail + dq.max - dq.head) % dq.max
}

func (dq *Queue) Show() {
	fmt.Printf("head=%d,tail=%d,size=%d\n", dq.head, dq.tail, dq.size())
}

func (dq *Queue) IsEmpty() bool {
	dq.RLock()
	defer dq.RUnlock()
	return dq.isEmpty()
}

func (dq *Queue) isEmpty() bool {
	return dq.tail == dq.head
}

func (dq *Queue) IsFull() bool {
	dq.RLock()
	defer dq.RUnlock()
	return dq.isFull()
}

func (dq *Queue) isFull() bool {
	return (dq.tail+1)%dq.max == dq.head
}
