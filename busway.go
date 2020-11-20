package busway

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Flag :nodoc:
type Flag int

const (
	defaultBufferCapacity = 1024
)

// These flag defines the flushing priorities
const (
	BNormal Flag = 1 << iota
	BImmediate
)

// Busway :nodoc:
type Busway interface {
	io.Writer
	io.Closer
	AddWriter(writer io.Writer)
	WriteRich(t Flag, b []byte) (n int, err error)
}

// Config :nodoc:
type Config struct {
	BufferCapacity int
}

// ConcreteBusway :nodoc:
type ConcreteBusway struct {
	cfg     Config
	writers atomic.Value

	lock     sync.Mutex
	closable chan bool

	messages chan []byte
}

// New :nodoc:
func New(cfg Config) Busway {
	if cfg.BufferCapacity == 0 {
		cfg.BufferCapacity = defaultBufferCapacity
	}

	cb := &ConcreteBusway{
		messages: make(chan []byte, cfg.BufferCapacity),
		lock:     sync.Mutex{},
	}

	cb.writers.Store([]io.Writer{})

	go func() {
		for msg := range cb.messages {
			cb.lock.Lock()
			cb.flush(msg)
			cb.lock.Unlock()

			if len(cb.messages) == 0 && cb.closable != nil {
				cb.closable <- true
			}
		}
	}()

	return cb
}

// AddWriter :nodoc:
func (cb *ConcreteBusway) AddWriter(writer io.Writer) {
	newWriters := append(cb.writers.Load().([]io.Writer), writer)
	cb.writers.Store(newWriters)
}

// WriteRich :nodoc:
func (cb *ConcreteBusway) WriteRich(tag Flag, b []byte) (n int, err error) {
	if cb.closable != nil {
		err = fmt.Errorf("cannot write: closing")
		return
	}

	switch true {
	case tag&BImmediate == 0:
		cb.lock.Lock()
		defer cb.lock.Unlock()
		if cb.flush(b); err != nil {
			return
		}

		break
	default:
		cb.messages <- b
	}

	return len(b), nil
}

// Write :nodoc:
func (cb *ConcreteBusway) Write(b []byte) (n int, err error) {
	return cb.WriteRich(BNormal, b)
}

// Close :nodoc:
func (cb *ConcreteBusway) Close() (err error) {
	cb.closable = make(chan bool)
	<-cb.closable
	return nil
}

func (cb *ConcreteBusway) flush(b []byte) (err error) {
	for _, writer := range cb.writers.Load().([]io.Writer) {
		var n int
		n, err = writer.Write(b)
		if err != nil {
			return
		}
		if n != len(b) {
			err = fmt.Errorf("incomplete data write")
			return
		}
	}
	return
}
