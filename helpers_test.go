package warnings_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/runbed/warnings"
)

// ExampleTap demonstrates how to use the Tap function to apply a side-effect to each warning.
func ExampleTap() {
	// create a new collector
	collector := warnings.NewCollector()
	defer collector.Close() // make sure to close the collector when done
	// attach the collector to a context
	ctx := warnings.Attach(context.Background(), collector)
	// use Tap to apply a side-effect to each warning
	ctx = warnings.Tap(ctx, func(wrr warnings.Warning) {
		fmt.Println("side-effect:", wrr.Warn())
	})
	// use Warn or Warnf to write warnings to the context
	warnings.Warnf(ctx, "this is a warning 1")
	warnings.Warnf(ctx, "this is a warning 2")
	// read all warnings from the collector
	wrrs, err := warnings.ReadAll(collector)
	if err != nil {
		// handle error
	}
	for _, wrr := range wrrs {
		fmt.Println(wrr.Warn())
	}
	// Output:
	// side-effect: this is a warning 1
	// side-effect: this is a warning 2
	// this is a warning 1
	// this is a warning 2
}

type multiWarn struct {
	details []string
}

func (w *multiWarn) Warn() string {
	return strings.Join(w.details, ", ")
}

// ExampleReduce demonstrates how to use the Reduce function to reduce warnings into a single value.
func ExampleReduce() {
	// create a new collector
	collector := warnings.NewCollector()
	defer collector.Close() // make sure to close the collector when done
	// attach the collector to a context
	ctx := warnings.Attach(context.Background(), collector)
	// use Reduce to reduce warnings into a single value
	ctx, flush := warnings.Reduce(ctx, func(acc *multiWarn, wrr warnings.Warning) *multiWarn {
		if acc == nil { // initialize the accumulator
			acc = new(multiWarn)
		}
		acc.details = append(acc.details, wrr.Warn())
		return acc
	})
	// use Warn or Warnf to write warnings to the context
	warnings.Warnf(ctx, "this is a warning 1")
	warnings.Warnf(ctx, "this is a warning 2")
	warnings.Warnf(ctx, "this is a warning 3")
	// flush the reduced warning
	flush()
	// read all warnings from the collector
	wrrs, err := warnings.ReadAll(collector)
	if err != nil {
		// handle error
	}
	for i, wrr := range wrrs {
		fmt.Printf("[%d]: %s\n", i, wrr.Warn())
	}
	// Output:
	// [0]: this is a warning 1, this is a warning 2, this is a warning 3
}

// ExampleFilter demonstrates how to use the Filter function to filter warnings.
func ExampleFilter() {
	// create a new collector
	collector := warnings.NewCollector()
	defer collector.Close() // make sure to close the collector when done
	// attach the collector to a context
	ctx := warnings.Attach(context.Background(), collector)
	// use Filter to filter warnings
	ctx = warnings.Filter(ctx, func(wrr warnings.Warning) bool {
		return !strings.HasPrefix(wrr.Warn(), "ignore")
	})
	// use Warn or Warnf to write warnings to the context
	warnings.Warnf(ctx, "this is a warning")
	warnings.Warnf(ctx, "ignore this warning")
	warnings.Warnf(ctx, "this is another warning")
	// read all warnings from the collector
	wrrs, err := warnings.ReadAll(collector)
	if err != nil {
		// handle error
	}
	for i, wrr := range wrrs {
		fmt.Printf("[%d]: %s\n", i, wrr.Warn())
	}
	// Output:
	// [0]: this is a warning
	// [1]: this is another warning
}

// ExampleMap demonstrates how to use the Map function to map warnings.
func ExampleMap() {
	// create a new collector
	collector := warnings.NewCollector()
	defer collector.Close() // make sure to close the collector when done
	// attach the collector to a context
	ctx := warnings.Attach(context.Background(), collector)
	// use Map to map warnings
	ctx = warnings.Map(ctx, func(wrr warnings.Warning) warnings.Warning {
		return warnings.New(strings.ToUpper(wrr.Warn()))
	})
	// use Warn or Warnf to write warnings to the context
	warnings.Warnf(ctx, "this is a warning")
	warnings.Warnf(ctx, "this is another warning")
	// read all warnings from the collector
	wrrs, err := warnings.ReadAll(collector)
	if err != nil {
		// handle error
	}
	for i, wrr := range wrrs {
		fmt.Printf("[%d]: %s\n", i, wrr.Warn())
	}
	// Output:
	// [0]: THIS IS A WARNING
	// [1]: THIS IS ANOTHER WARNING
}

func TestMap(t *testing.T) {
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	ctx = warnings.Map(ctx, func(wrr warnings.Warning) warnings.Warning {
		return warnings.New(strings.ToUpper(wrr.Warn()))
	})
	warnings.Warn(ctx, warnings.New("test"))
	if len(w.buf) != 1 {
		t.Fatalf("expected 1 warning, got %v", len(w.buf))
	}
	if got, want := w.buf[0].Warn(), "TEST"; got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestMapNoWriter(t *testing.T) {
	ctx := warnings.Map(context.Background(), func(wrr warnings.Warning) warnings.Warning {
		return warnings.New(strings.ToUpper(wrr.Warn()))
	})
	err := warnings.Warn(ctx, warnings.New("test"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestFilter(t *testing.T) {
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	ctx = warnings.Filter(ctx, func(wrr warnings.Warning) bool {
		return wrr.Warn() != "ignore"
	})
	warnings.Warn(ctx, warnings.New("this"))
	warnings.Warn(ctx, warnings.New("ignore"))
	warnings.Warn(ctx, warnings.New("that"))
	if len(w.buf) != 2 {
		t.Fatalf("expected 2 warnings, got %v", len(w.buf))
	}
	if got := w.buf[0].Warn(); got != "this" {
		t.Fatalf("expected this, got %v", got)
	}
	if got := w.buf[1].Warn(); got != "that" {
		t.Fatalf("expected that, got %v", got)
	}
}

func TestFilterNoWriter(t *testing.T) {
	ctx := warnings.Filter(context.Background(), func(wrr warnings.Warning) bool {
		return wrr.Warn() != "ignore"
	})
	err := warnings.Warn(ctx, warnings.New("ignore"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestReduce(t *testing.T) {
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	ctx, flush := warnings.Reduce(ctx, func(acc *multiWarn, wrr warnings.Warning) *multiWarn {
		if acc == nil {
			acc = new(multiWarn)
		}
		acc.details = append(acc.details, wrr.Warn())
		return acc
	})
	warnings.Warn(ctx, warnings.New("this"))
	warnings.Warn(ctx, warnings.New("that"))
	flush()
	if len(w.buf) != 1 {
		t.Fatalf("expected 1 warning, got %v", len(w.buf))
	}
	if got := w.buf[0].Warn(); got != "this, that" {
		t.Fatalf("expected this, that, got %v", got)
	}
}

func TestReduceNoWriter(t *testing.T) {
	ctx, flush := warnings.Reduce(context.Background(), func(acc *multiWarn, wrr warnings.Warning) *multiWarn {
		if acc == nil {
			acc = new(multiWarn)
		}
		acc.details = append(acc.details, wrr.Warn())
		return acc
	})
	warnings.Warn(ctx, warnings.New("this"))
	warnings.Warn(ctx, warnings.New("that"))
	flush()
}

func TestReduceNoWarning(t *testing.T) {
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	_, flush := warnings.Reduce(ctx, func(acc *multiWarn, wrr warnings.Warning) *multiWarn {
		if acc == nil {
			acc = new(multiWarn)
		}
		acc.details = append(acc.details, wrr.Warn())
		return acc
	})
	flush()
	if len(w.buf) > 0 {
		t.Fatalf("expected no warnings, got %v", len(w.buf))
	}
}

func TestTap(t *testing.T) {
	touched := false
	w := &mockWriter{}
	ctx := warnings.Attach(context.Background(), w)
	ctx = warnings.Tap(ctx, func(wrr warnings.Warning) {
		touched = true
	})
	warnings.Warn(ctx, warnings.New("test"))
	if !touched {
		t.Fatalf("expected touched, got not touched")
	}
	if len(w.buf) != 1 {
		t.Fatalf("expected 1 warning, got %v", len(w.buf))
	}
}

func TestTapNoWriter(t *testing.T) {
	touched := false
	ctx := warnings.Tap(context.Background(), func(wrr warnings.Warning) {
		touched = true
	})
	warnings.Warn(ctx, warnings.New("test"))
	if touched {
		t.Fatalf("expected not touched, got touched")
	}
}
