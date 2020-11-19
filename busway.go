package busway

import (
	"fmt"
	"io"
	"sync/atomic"
)

const (
	// How many log messages can be buffered until the call blocks
	bufferCapacity = 1024
)

// Busway :nodoc:
type Busway interface {
	io.Writer
	io.Closer
	AddWriter(writer io.Writer)
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

	lobuf chan []byte
	hibuf chan []byte
}

// New creates a new Busway.
func New(cfg Config) Busway {
	log := &ConcreteBusway{
		lobuf: make(chan []byte, bufferCapacity),
	}

	log.writers.Store([]io.Writer{})

	go func() {
		for msg := range log.lobuf {
			log.write(msg)
		}
	}()

	return log
}

// AddWriter adds an output to the log.
func (log *ConcreteBusway) AddWriter(writer io.Writer) {
	newWriters := append(log.writers.Load().([]io.Writer), writer)
	log.writers.Store(newWriters)
}

// Info writes non-critical information to the log.
// Unlike Error, it does not guarantee that the message will have been
// written persistenly to disk at the time this function returns.
func (log *ConcreteBusway) Info(format string, values ...interface{}) {
	fmt.Fprintf(log, format+"\n", values...)
}

// Error writes critical information to the log.
// It will instantly flush the I/O buffers and guarantees that the message
// will have been written persistenly to disk at the time this function returns.
func (log *ConcreteBusway) Error(format string, values ...interface{}) {
	fmt.Fprintf(log, format+"\n", values...)
	// TODO: Flush.
}

// Write implements the io.Writer interface.
// As long as buffer capacity is available,
// this call will not block and have O(1) behaviour,
// regardless of how many writers are used.
func (log *ConcreteBusway) Write(b []byte) (n int, err error) {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	log.lobuf <- tmp
	return len(b), nil
}

// Close implements the io.Closer interface.
func (log *ConcreteBusway) Close() (err error) {
	return nil
}

// write writes the given slice of bytes to all registered writers immediately.
func (log *ConcreteBusway) write(b []byte) {
	for _, writer := range log.writers.Load().([]io.Writer) {
		_, _ = writer.Write(b)
	}
}
