package warnings

import "errors"

// Writer is the interface for writing warnings.
type Writer interface {
	// WriteWarning writes the warning.
	WriteWarning(wrr Warning) error
}

type multiWriter struct {
	writers []Writer
}

// NewMultiWriter returns a Writer that duplicates its writes to all the provided writers.
func NewMultiWriter(writers ...Writer) Writer {
	return &multiWriter{writers}
}

func (w *multiWriter) WriteWarning(wrr Warning) error {
	var errs []error
	for _, w := range w.writers {
		if err := w.WriteWarning(wrr); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
