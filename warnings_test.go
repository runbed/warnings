package warnings_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/runbed/warnings"
)

// Example demonstrates how t o use the warnings package to read and write warnings.
func Example() {
	// create a new collector
	collector := warnings.NewCollector()
	defer collector.Close() // make sure to close the collector when done
	// attach the collector to a context
	ctx := warnings.Attach(context.Background(), collector)
	// use Warn or Warnf to write warnings to the context
	warnings.Warnf(ctx, "this is a warning 1")
	warnings.Warnf(ctx, "this is a warning 2")
	// use Scanner to read warnings one by one
	scanner := warnings.NewScanner(collector)
	for scanner.Scan() {
		wrr := scanner.Warning()
		fmt.Println(wrr.Warn())
	}
	if err := scanner.Err(); err != nil {
		// handle error
	}
	// Output:
	// this is a warning 1
	// this is a warning 2
}

func TestNew(t *testing.T) {
	want := "test-warning"
	wrr := warnings.New(want)
	if got := wrr.Warn(); got != want {
		t.Errorf("expected %s, got %v", want, got)
	}
	if got, ok := wrr.(fmt.Stringer); !ok {
		t.Errorf("expected warning to implement fmt.Stringer")
	} else if got := got.String(); got != want {
		t.Errorf("expected %s, got %v", want, got)
	}
	if got, ok := wrr.(json.Marshaler); !ok {
		t.Errorf("expected warning to implement json.Marshaler")
	} else {
		got, err := got.MarshalJSON()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !bytes.Equal(got, []byte(`"`+want+`"`)) {
			t.Errorf("expected %v, got %v", want, got)
		}
	}
}

func TestWarn(t *testing.T) {
	want := warnings.New("test-warning")
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	err := warnings.Warn(ctx, want)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(w.buf) != 1 || w.buf[0] != want {
		t.Fatalf("expected %v, got %v", want, w.buf)
	}
}

func TestWarnNoWriter(t *testing.T) {
	want := warnings.New("test-warning")
	ctx := context.Background()
	err := warnings.Warn(ctx, want)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestWarnError(t *testing.T) {
	wantErr := fmt.Errorf("test-error")
	w := &mockWriter{result: wantErr}
	ctx := warnings.Attach(context.Background(), w)
	wantWrr := warnings.New("test-warning")
	err := warnings.Warn(ctx, wantWrr)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if len(w.buf) > 0 {
		t.Fatalf("expected no warnings, got %v", w.buf)
	}
}

func TestWarnf(t *testing.T) {
	w := &mockWriter{}
	want := "test-warning: sub-warning"
	ctx := warnings.Attach(context.Background(), w)
	err := warnings.Warnf(ctx, "test-warning: %s", warnings.New("sub-warning"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(w.buf) != 1 || w.buf[0].Warn() != want {
		t.Fatalf("expected test-1, got %v", w.buf)
	}
}

func TestAttach(t *testing.T) {
	wrrs := []warnings.Warning{
		warnings.New("test-1"),
		warnings.New("test-2"),
	}
	writers := make([]*mockWriter, len(wrrs))
	for i := range writers {
		writers[i] = new(mockWriter)
	}
	ctx := warnings.Attach(context.Background(), writers[0])
	warnings.Warn(ctx, wrrs[0])
	ctx = warnings.Attach(ctx, writers[1])
	warnings.Warn(ctx, wrrs[1])
	if len(writers[0].buf) != 2 {
		t.Fatalf("expected 2 warnings, got %v", writers[0].buf)
	}
	if len(writers[1].buf) != 1 {
		t.Fatalf("expected 1 warning, got %v", writers[1].buf)
	}
}

func TestDetach(t *testing.T) {
	wantWrr := warnings.New("test-warning")
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	ctx = warnings.Detach(ctx)
	err := warnings.Warn(ctx, wantWrr)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(w.buf) > 0 {
		t.Fatalf("expected no warnings, got %v", w.buf)
	}
}

func TestDetachNoWriter(t *testing.T) {
	ctx := context.Background()
	ctx = warnings.Detach(ctx)
	err := warnings.Warn(ctx, warnings.New("test-warning"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
