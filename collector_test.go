package warnings_test

import (
	"testing"

	"github.com/runbed/warnings"
)

func TestCollector_Close(t *testing.T) {
	c := warnings.NewCollector()
	err := c.Close()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	err = c.Close()
	if err != warnings.ErrClosed {
		t.Fatalf("expected %v, got %v", warnings.ErrClosed, err)
	}
}

func TestCollector_WriteWarning(t *testing.T) {
	c := warnings.NewCollector()
	err := c.WriteWarning(warnings.New("test-1"))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	err = c.Close()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	err = c.WriteWarning(warnings.New("test-2"))
	if err != warnings.ErrClosed {
		t.Fatalf("expected %v, got %v", warnings.ErrClosed, err)
	}
}

func TestCollector_ReadWarning(t *testing.T) {
	c := warnings.NewCollector()
	err := c.WriteWarning(warnings.New("test-1"))
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	w, err := c.ReadWarning()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if w.Warn() != "test-1" {
		t.Fatalf("expected test-1, got %v", w.Warn())
	}
	err = c.Close()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	w, err = c.ReadWarning()
	if err != warnings.ErrClosed {
		t.Fatalf("expected %v, got %v", warnings.ErrClosed, err)
	}
	if w != nil {
		t.Fatalf("expected nil, got %v", w)
	}
}
