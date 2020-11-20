package wrcrux

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// Mode :nodoc:
type Mode int

const (
	defaultBufferCapacity = 1024
)

// These flag defines the flushing priorities
const (
	Buffered Mode = 1 << iota
	Immediate
)

// Wrcrux :nodoc:
type Wrcrux interface {
	io.Writer
	io.Closer
	AddWriter(writer io.Writer)
	XWrite(t Mode, b []byte) (n int, err error)
}

// Config :nodoc:
type Config struct {
	BufferCapacity int
	DefaultMode    Mode
}

// ConcreteWrcrux :nodoc:
type ConcreteWrcrux struct {
	cfg     Config
	writers atomic.Value

	lock     sync.Mutex
	closable chan bool

	messages chan []byte
}

// New :nodoc:
func New(cfg Config) Wrcrux {
	if cfg.BufferCapacity == 0 {
		cfg.BufferCapacity = defaultBufferCapacity
	}
	if cfg.DefaultMode == 0 {
		cfg.DefaultMode = Buffered
	}

	cb := &ConcreteWrcrux{
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
func (cb *ConcreteWrcrux) AddWriter(writer io.Writer) {
	newWriters := append(cb.writers.Load().([]io.Writer), writer)
	cb.writers.Store(newWriters)
}

// XWrite :nodoc:
func (cb *ConcreteWrcrux) XWrite(tag Mode, b []byte) (n int, err error) {
	if cb.closable != nil {
		err = fmt.Errorf("cannot write: closing")
		return
	}

	switch true {
	case tag == Immediate:
		cb.lock.Lock()
		if err = cb.flush(b); err != nil {
			cb.lock.Unlock()
			return
		}
		cb.lock.Unlock()

		break
	default:
		cb.messages <- b
	}

	return len(b), nil
}

// Write :nodoc:
func (cb *ConcreteWrcrux) Write(b []byte) (n int, err error) {
	return cb.XWrite(Buffered, b)
}

// Close :nodoc:
func (cb *ConcreteWrcrux) Close() (err error) {
	cb.closable = make(chan bool)
	<-cb.closable
	return nil
}

func (cb *ConcreteWrcrux) flush(b []byte) (err error) {
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
