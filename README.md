# Warnings
[![License](https://img.shields.io/badge/license-mit-green.svg)](https://github.com/runbed/warnings/blob/main/LICENSE)
[![Go](https://github.com/runbed/warnings/actions/workflows/go.yml/badge.svg)](https://github.com/runbed/warnings/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/runbed/warnings/badge.svg?branch=main)](https://coveralls.io/github/runbed/warnings?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/runbed/warnings)](https://goreportcard.com/report/github.com/runbed/warnings)
[![GoDoc](https://godoc.org/github.com/runbed/warnings?status.svg)](http://godoc.org/github.com/runbed/warnings)
[![Release](https://img.shields.io/github/release/runbed/warnings.svg)](https://github.com/runbed/warnings/releases/latest)

`warnings` package provides mechanisms for capturing diagnostics using `context.Context`.
It allows to easily capture warnings without modifying existing function signatures.
Warnings are captured using a `Collector`, which is attached to the `context.Context`.

## Overview

The package offers the following functionalities:

- **Capturing**: The package provides mechanisms to collect warnings generated during the execution of a Go program. Warnings can be captured and stored for later analysis or reporting.
- **Filtering**: Developers can apply filters to selectively capture or ignore certain types of warnings based on specified criteria. This allows for targeted handling of warnings that meet specific conditions.
- **Transformation of Warnings**: Warnings can be transformed or modified in various ways to suit different requirements. This functionality enables developers to manipulate warning messages, format them differently, or perform other operations before storing or handling them.
- **Side-Effects**: The package supports the application of side-effects to warnings, such as logging or triggering additional actions based on the occurrence of specific warnings. This capability allows for flexible handling of warnings beyond simple collection and storage.

## Installation

To use this package in your Go project, you can import it using:

```go
import "github.com/runbed/warnings"
```

## Usage

To start capturing warnings create a new `Collector`:

```go
collector := warnings.NewCollector()
defer collector.Close() // close before exiting
```

Attach this collector to your context:

```go
ctx := warnings.Attach(context.Background(), collector)
```

Then, you can write warnings to the context:

```go
warnings.Warnf(ctx, "this is a warning")
```

Finally, capture all of warnings back from the collector:

```go
wrrs, err := warnings.ReadAll(collector)
```
### Helpers

#### Filter

This example demonstrates how to use the `Filter` function to filter warnings:

```go
// filter warnings
ctx = warnings.Filter(ctx, func(wrr warnings.Warning) bool {
    return !strings.HasPrefix(wrr.Warn(), "ignore:")
})
// This warning will be captured
warnings.Warnf(ctx, "capture: warning 1")
// This warning will be ignored
warnings.Warnf(ctx, "ignore: warning 2") 
```

#### Map

Transform each written warning.

```go
ctx = warnings.Map(ctx, func(wrr warnings.Warning) warnings.Warning {
    return warnings.New(strings.ToUpper(wrr.Warn()))
})
// This warning will be transformed to "WARNING 1"
warnings.Warnf(ctx, "warning 1")
```

#### Reduce

Combine all written warnings into a single warning.

```go
// create a custom warning
type multiWarn struct {
	details []string
}

func (w *multiWarn) Warn() string {
	return strings.Join(w.details, ", ")
}

// reduce warnings
ctx, flush := warnings.Reduce(ctx, func(acc *multiWarn*, wrr warnings.Warning) *multiWarn {
    if acc == nil {
        acc = new(multiWarn)
    }
    acc.details = append(acc.details, wrr.Warn())
    return acc
})
// flush on exit
defer flush()
// Write warnings
warnings.Warnf(ctx, "warning 1")
warnings.Warnf(ctx, "warning 2")
warnings.Warnf(ctx, "warning 3")
// The captured warning will be:
// &multiWarn{"warning 1", "warning 2", "warning 3"}
```

#### Tap

It does not modify the warnings or the context but is useful for side effects like logging.

```go
ctx = warnings.Tap(ctx, func(wrr warnings.Warning) {
    slog.Warn(wrr.Warn())
})
// Now every new warning will be logged using `slog`
warnings.Warnf(ctx, "this is a warning")
warnings.Warnf(ctx, "this is another warning")
```

## Contributing

Thank you for your interest in contributing to the `warnings` Go library! We welcome and appreciate any contributions, whether they be bug reports, feature requests, or code changes.

If you've found a bug, please create an issue in the GitHub repository describing the problem, including any relevant error messages and a minimal reproduction of the issue.

## License

`warnings` is licensed under the [MIT License](LICENSE).