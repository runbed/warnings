package warnings

import (
	"io"
	"sync"
)

// Collector type is used to capture warnings.
// It implements the [Reader], [Writer] and [io.Closer] interfaces.
// Read operation are non-blocking and returns [io.EOF] when there are no more warnings in the buffer.
// The collector is thread-safe. It is safe to read and write warnings concurrently.
type Collector struct {
	buf    []Warning
	mtx    sync.Mutex
	closed bool
}

// NewCollector returns a new Collector.
func NewCollector() *Collector {
	return new(Collector)
}

// Close closes the collector.
func (c *Collector) Close() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.closed {
		return ErrClosed
	}
	c.closed = true
	c.buf = nil
	return nil
}

// WriteWarning writes a warning to the collector.
func (c *Collector) WriteWarning(wrr Warning) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.closed {
		return ErrClosed
	}
	c.buf = append(c.buf, wrr)
	return nil
}

// ReadWarning reads a warning from the collector.
func (c *Collector) ReadWarning() (Warning, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.closed {
		return nil, ErrClosed
	}
	if len(c.buf) == 0 {
		return nil, io.EOF
	}
	w := c.buf[0]
	c.buf = c.buf[1:]
	return w, nil
}
