// Package warnings implements mechanisms for capturing diagnostics using [context.Context].
//
// It is designed to provide an easy way to capture warnings without modifying existing function signature.
//
// To start capturing warnings, you need to attach a [Collector] to the context using [Attach] function.
// This will create a new context with the collector attached and all the warnings written to
// the context will be captured.
//
//	collector := warnings.NewCollector()
//	defer collector.Close()
//	ctx := warnings.Attach(context.Background(), collector)
//
// Use [Warn] or [Warnf] functions to write warnings to the context. H
//
//	warnings.Warnf(ctx, "this is a warning")
//
// To read all the warnings from the collector, use [ReadAll] function
//
//	wrrs, err := warnings.ReadAll(collector)
//
// Or you can use [Scanner] function to read warnings one by one.
//
//	scanner := warnings.NewScanner(collector)
//	for scanner.Scan() {
//		wrr := scanner.Warning()
//	}
//	if err := scanner.Err(); err != nil {
//		// handle error
//	}
//
// If you need a new context that does not collect warnings anymore, use [Detach] function.
//
//	ctx = warnings.Detach(ctx)
//
// Use [Map], [Filter], [Reduce] or [Tap] helper functions to apply transformations,
// filters or side-effects to the warnings.
package warnings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Warning is an interface representing a warning.
type Warning interface {
	Warn() string
}

type warningString struct {
	s string
}

func (wrr *warningString) Warn() string {
	return wrr.s
}

func (wrr *warningString) String() string {
	return wrr.s
}

func (wrr *warningString) MarshalJSON() ([]byte, error) {
	return json.Marshal(wrr.s)
}

// New creates a new warning from a given string.
func New(str string) Warning {
	return &warningString{str}
}

type writerKey struct{}

func setWriter(ctx context.Context, w Writer) context.Context {
	return context.WithValue(ctx, writerKey{}, w)
}

func getWriter(ctx context.Context) Writer {
	w, ok := ctx.Value(writerKey{}).(Writer)
	if !ok {
		return nil
	}
	return w
}

func resetWriter(ctx context.Context, w Writer) context.Context {
	return setWriter(setWriter(ctx, nil), w)
}

// Warn writes warnings to the context. When multiple warnings are provided, they are written in order.
// If no writer is attached to the context, it does nothing and returns nil.
// If any of the warnings fail to write, all the warnings are returned as one error.
func Warn(ctx context.Context, wrrs ...Warning) error {
	w := getWriter(ctx)
	if w == nil {
		return nil
	}
	var errs []error
	for _, wrr := range wrrs {
		if err := w.WriteWarning(wrr); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Warnf is a helper function that formats the warning and writes it to the context.
// If the format string contains any [Warning] arguments, they are converted to strings before formatting.
func Warnf(ctx context.Context, format string, args ...any) error {
	for i, arg := range args {
		if wrr, ok := arg.(Warning); ok {
			args[i] = wrr.Warn()
		}
	}
	return Warn(ctx, &warningString{fmt.Sprintf(format, args...)})
}

// Attach returns a new context that collects warnings using the provided writer.
// If a writer is already attached to the context, it creates a new writer that writes to both.
func Attach(ctx context.Context, w Writer) context.Context {
	if found := getWriter(ctx); found != nil {
		w = NewMultiWriter(found, w)
	}
	return setWriter(ctx, w)
}

// Detach returns a new context that does not propagate warnings up the chain.
func Detach(ctx context.Context) context.Context {
	w := getWriter(ctx)
	if w == nil {
		return ctx
	}
	return setWriter(ctx, nil)
}
