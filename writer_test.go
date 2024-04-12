package warnings_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/runbed/warnings"
)

type mockWriter struct {
	buf    []warnings.Warning
	result error
}

func (w *mockWriter) WriteWarning(wrr warnings.Warning) error {
	if w.result != nil {
		return w.result
	}
	w.buf = append(w.buf, wrr)
	return w.result
}

func TestNewMultiWriter(t *testing.T) {
	writers := []*mockWriter{
		{result: nil},
		{result: nil},
	}
	w := warnings.NewMultiWriter(writers[0], writers[1])
	wantWrr := warnings.New("test")
	err := w.WriteWarning(wantWrr)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	for _, writer := range writers {
		if len(writer.buf) != 1 || writer.buf[0] != wantWrr {
			t.Fatalf("expected %v, got %v", wantWrr, writer.buf)
		}
	}
}

func TestNewMultiWriterErrors(t *testing.T) {
	writers := []*mockWriter{
		{result: nil},
		{result: fmt.Errorf("test-error-1")},
		{result: fmt.Errorf("test-error-2")},
	}
	w := warnings.NewMultiWriter(writers[0], writers[1], writers[2])
	wantWrr := warnings.New("test")
	err := w.WriteWarning(wantWrr)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, writers[1].result) {
		t.Fatalf("expected %v, got %v", writers[1].result, err)
	}
	if !errors.Is(err, writers[2].result) {
		t.Fatalf("expected %v, got %v", writers[2].result, err)
	}
}
