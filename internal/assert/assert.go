package assert

import (
	"strings"
	"testing"
)

func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("assert.True: got false")
	}
}

func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("assert.False: got true")
	}
}

func Equal[T comparable](t *testing.T, got, expected T) {
	t.Helper()
	if got != expected {
		t.Errorf("assert.Equal: got %v, expected %v", got, expected)
	}
}

func NotEqual[T comparable](t *testing.T, got, illegal T) {
	t.Helper()
	if got == illegal {
		t.Errorf("assert.NotEqual: got %v", got)
	}
}

func Nil(t *testing.T, got any) {
	t.Helper()
	if got != nil {
		t.Fatalf("assert.Nil: got %v", got)
	}
}

func NotNil(t *testing.T, got any) {
	t.Helper()
	if got == nil {
		t.Fatalf("assert.NotNil: got nil")
	}
}

func Contains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("assert.Contains: %q did not contain %q", s, substr)
	}
}

func HasPrefix(t *testing.T, s, prefix string) {
	t.Helper()
	if !strings.HasPrefix(s, prefix) {
		t.Errorf("assert.HasPrefix: %q did not have prefix %q", s, prefix)
	}
}

func HasSuffix(t *testing.T, s, suffix string) {
	t.Helper()
	if !strings.HasSuffix(s, suffix) {
		t.Errorf("assert.HasSuffix: %q did not have suffix %q", s, suffix)
	}
}
