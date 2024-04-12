package warnings

import (
	"errors"
	"io"
)

// Scanner reads warnings from an underlying reader and provides a way to iterate over them.
// The Scanner is not thread-safe. It is the caller's responsibility to ensure that
// the Scanner is not being used concurrently.
type Scanner interface {
	// Scan advances the scanner to the next warning.
	Scan() bool
	// Warning returns the current warning.
	Warning() Warning
	// Err returns the first non-EOF error that was encountered by the scanner.
	Err() error
}

// NewScanner returns a new Scanner.
func NewScanner(r Reader) Scanner {
	return &scanner{r, nil, nil}
}

type scanner struct {
	r   Reader
	wrr Warning
	err error
}

func (s *scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	wrr, err := s.r.ReadWarning()
	if errors.Is(err, io.EOF) {
		s.wrr = nil
		return false
	} else if err != nil {
		s.err = err
		return false
	}
	s.wrr = wrr
	return true
}

func (s *scanner) Warning() Warning {
	return s.wrr
}

func (s *scanner) Err() error {
	return s.err
}
