package warnings_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/runbed/warnings"
)

func TestScanner(t *testing.T) {
	want := []warnings.Warning{
		warnings.New("test-1"),
		warnings.New("test-2"),
	}
	r := &mockReader{
		[]mockReaderResult{
			{want[0], nil},
			{want[1], nil},
			{nil, io.EOF},
		},
	}
	scanner := warnings.NewScanner(r)
	for i, wrr := range want {
		if !scanner.Scan() {
			t.Fatalf("expected to scan warning %v", i)
		}
		if got := scanner.Warning(); got != wrr {
			t.Errorf("expected %v, got %v", wrr, got)
		}
	}
	if scanner.Scan() {
		t.Fatalf("expected to not scan any more warnings")
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestScanner_EOF(t *testing.T) {
	r := &mockReader{
		[]mockReaderResult{
			{warnings.New("test-1"), nil},
			{warnings.New("test-2"), nil},
			{nil, io.EOF},
		},
	}
	scanner := warnings.NewScanner(r)
	if !scanner.Scan() {
		t.Fatalf("expected to scan warning 0")
	}
	if !scanner.Scan() {
		t.Fatalf("expected to scan warning 1")
	}
	if scanner.Scan() {
		t.Fatalf("expected to not scan any more warnings")
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestScanner_UnexpectedError(t *testing.T) {
	wantErr := fmt.Errorf("test-error")
	r := &mockReader{
		[]mockReaderResult{
			{warnings.New("test-1"), nil},
			{warnings.New("test-2"), nil},
			{nil, wantErr},
		},
	}
	scanner := warnings.NewScanner(r)
	if !scanner.Scan() {
		t.Fatalf("expected to scan warning 0")
	}
	if !scanner.Scan() {
		t.Fatalf("expected to scan warning 1")
	}
	if scanner.Scan() {
		t.Fatalf("expected to not scan any more warnings")
	}
	if err := scanner.Err(); err != wantErr {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if scanner.Scan() {
		t.Fatalf("expected to not scan any more warnings")
	}
	if err := scanner.Err(); err != wantErr {
		t.Fatalf("expected same error %v, got %v", wantErr, err)
	}
}

func TestScanner_Empty(t *testing.T) {
	r := &mockReader{
		[]mockReaderResult{
			{nil, io.EOF},
		},
	}
	scanner := warnings.NewScanner(r)
	if scanner.Scan() {
		t.Fatalf("expected to not scan any warnings")
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
