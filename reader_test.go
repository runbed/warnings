package warnings_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/runbed/warnings"
)

type mockReader struct {
	results []mockReaderResult
}

func (m *mockReader) ReadWarning() (warnings.Warning, error) {
	if len(m.results) == 0 {
		return nil, fmt.Errorf("no more results")
	}
	result := m.results[0]
	m.results = m.results[1:]
	return result.wrr, result.err
}

type mockReaderResult struct {
	wrr warnings.Warning
	err error
}

func TestReadAll(t *testing.T) {
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
	wrrs, err := warnings.ReadAll(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l := len(wrrs); l != 2 {
		t.Fatalf("expected 2 warnings, got %v", l)
	}
	for i, wrr := range want {
		if got := wrrs[i]; got != wrr {
			t.Errorf("expected %v, got %v", wrr, got)
		}
	}
}

func TestReadAll_UnexpectedError(t *testing.T) {
	wantErr := fmt.Errorf("test-error")
	r := &mockReader{
		[]mockReaderResult{
			{warnings.New("test-1"), nil},
			{warnings.New("test-2"), nil},
			{nil, wantErr},
		},
	}
	wrrs, err := warnings.ReadAll(r)
	if err != wantErr {
		t.Errorf("expected error %v, got %v", wantErr, err)
	}
	if len(wrrs) > 0 {
		t.Errorf("expected no warnings, got %v", wrrs)
	}
}

func TestReadAll_Empty(t *testing.T) {
	r := &mockReader{
		[]mockReaderResult{
			{nil, io.EOF},
		},
	}
	wrrs, err := warnings.ReadAll(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if l := len(wrrs); l > 0 {
		t.Fatalf("expected no warnings, got %v", l)
	}
}
