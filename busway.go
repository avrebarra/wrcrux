package busway

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type Tag int

const (
	// How many log messages can be buffered until the call blocks
	bufferCapacity = 1024
)

const (
	TNone Tag = iota
	TImmedate
)

// Busway :nodoc:
type Busway interface {
	io.Writer
	io.Closer
	AddWriter(writer io.Writer)
	WriteRich(t Tag, b []byte) (n int, err error)
}

// Config :nodoc:
type Config struct {
}

// ConcreteBusway is a log data source used for a specific part of your application,
// e.g. "web", "database", "api" or other categories. It can be connected
// to multiple writers.
type ConcreteBusway struct {
	cfg     Config
	writers atomic.Value

	lock     sync.Mutex
	closable chan bool

	messages chan []byte
}

// New creates a new Busway.
func New(cfg Config) Busway {
	cb := &ConcreteBusway{
		messages: make(chan []byte, bufferCapacity),
		lock:     sync.Mutex{},
	}

	cb.writers.Store([]io.Writer{})

	go func() {
		for msg := range cb.messages {
			cb.lock.Lock()
			cb.write(msg)
			cb.lock.Unlock()

			if len(cb.messages) == 0 && cb.closable != nil {
				cb.closable <- true
			}
		}
	}()

	return cb
}

// AddWriter adds an output to the cb.
func (cb *ConcreteBusway) AddWriter(writer io.Writer) {
	newWriters := append(cb.writers.Load().([]io.Writer), writer)
	cb.writers.Store(newWriters)
}

// WriteRich is writer with richer interface to control writing behavior
func (cb *ConcreteBusway) WriteRich(tag Tag, b []byte) (n int, err error) {
	if cb.closable != nil {
		err = fmt.Errorf("cannot write: closing")
		return
	}

	switch true {
	case tag == TImmedate:
		cb.lock.Lock()
		cb.write(b)
		cb.lock.Unlock()
		break
	default:
		cb.messages <- b
	}

	return len(b), nil
}

// Write implements the io.Writer interface.
// As long as buffer capacity is available,
// this call will not block and have O(1) behaviour,
// regardless of how many writers are used.
func (cb *ConcreteBusway) Write(b []byte) (n int, err error) {
	cb.WriteRich(0, b)
	return len(b), nil
}

// Close implements the io.Closer interface.
func (cb *ConcreteBusway) Close() (err error) {
	cb.closable = make(chan bool)
	<-cb.closable
	return nil
}

// write writes the given slice of bytes to all registered writers immediately.
func (cb *ConcreteBusway) write(b []byte) (err error) {
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
