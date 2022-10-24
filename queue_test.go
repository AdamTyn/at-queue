package queue

import (
	"errors"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	capacity := 3
	q := NewQueue(capacity, time.Second)
	_ = q.Put(func() error {
		return errors.New("1 put")
	})
	_ = q.Put(func() error {
		return errors.New("2 put")
	})
	_ = q.Put(func() error {
		return errors.New("3 put")
	})
	_ = q.Put(func() error {
		return errors.New("4 put")
	})
	q.Show()
	q.Run()
	defer q.Close()
}
