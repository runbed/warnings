package warnings

import "fmt"

var (
	// ErrClosed is returned when the warning stream is closed.
	ErrClosed = fmt.Errorf("warning stream is closed")
)
