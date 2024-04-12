package warnings

import (
	"errors"
	"io"
)

// Reader is the interface that wraps the basic ReadWarning method.
type Reader interface {
	// ReadWarning reads one warning from the reader.
	// If there are no more warnings, it returns [io.EOF].
	// If the reader is closed, it returns [ErrClosed].
	ReadWarning() (Warning, error)
}

// ReadAll reads all the warnings from the reader.
// It stops reading when it encounters an error or [io.EOF].
func ReadAll(r Reader) ([]Warning, error) {
	var result []Warning
	for {
		w, err := r.ReadWarning()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}
		result = append(result, w)
	}
	return result, nil
}
