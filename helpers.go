package warnings

import (
	"context"
)

// Map returns a new context that transforms each written warning using the provided function.
func Map(ctx context.Context, fn func(wrr Warning) Warning) context.Context {
	w := getWriter(ctx)
	if w == nil {
		return ctx
	}
	ctx = resetWriter(ctx, &mapWriter{w: w, fn: fn})
	return ctx
}

type mapWriter struct {
	w  Writer
	fn func(wrr Warning) Warning
}

func (mw *mapWriter) WriteWarning(wrr Warning) error {
	return mw.w.WriteWarning(mw.fn(wrr))
}

// Filter returns a new context that filters written warnings using the provided function.
func Filter(ctx context.Context, fn func(wrr Warning) bool) context.Context {
	w := getWriter(ctx)
	if w == nil {
		return ctx
	}
	ctx = resetWriter(ctx, &filterWriter{w: w, fn: fn})
	return ctx
}

type filterWriter struct {
	w  Writer
	fn func(wrr Warning) bool
}

func (fw *filterWriter) WriteWarning(wrr Warning) error {
	if fw.fn(wrr) {
		return fw.w.WriteWarning(wrr)
	}
	return nil
}

// Reduce returns a new context that reduces written warnings using the provided function.
// It also returns a flush() function that once called, writes the reduced warning to the underlying writer.
// If no warnings are written, it does nothing.
func Reduce[T Warning](ctx context.Context, fn func(acc T, wrr Warning) T) (_ context.Context, flush func()) {
	w := getWriter(ctx)
	if w == nil {
		return ctx, func() {}
	}
	input := NewCollector()
	ctx = resetWriter(ctx, input)
	return ctx, func() {
		defer input.Close()
		acc := *new(T)
		wrrs, err := ReadAll(input)
		if err != nil || len(wrrs) == 0 {
			return
		}
		for _, wrr := range wrrs {
			acc = fn(acc, wrr)
		}
		_ = w.WriteWarning(acc)
	}
}

// Tap returns a new context that taps written warnings using the provided function.
// It does not modify the warnings or the context but is useful for side effects like logging.
func Tap(ctx context.Context, fn func(wrr Warning)) context.Context {
	w := getWriter(ctx)
	if w == nil {
		return ctx
	}
	ctx = resetWriter(ctx, &tapWriter{w: w, fn: fn})
	return ctx
}

type tapWriter struct {
	w  Writer
	fn func(wrr Warning)
}

func (tw *tapWriter) WriteWarning(wrr Warning) error {
	tw.fn(wrr)
	return tw.w.WriteWarning(wrr)
}
