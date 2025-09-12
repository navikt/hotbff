package assert

import "testing"

func True(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("assert.True: expected true, got false")
	}
}

func False(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("assert.False: expected false, got true")
	}
}

func Equal[T comparable](t *testing.T, got, expected T) {
	t.Helper()
	if got != expected {
		t.Errorf("assert.Equal: expected %v, got %v", expected, got)
	}
}

func NotEqual[T comparable](t *testing.T, got, unexpected T) {
	t.Helper()
	if got == unexpected {
		t.Errorf("assert.NotEqual: got unexpected value %v", unexpected)
	}
}

func Nil(t *testing.T, got any) {
	t.Helper()
	if got != nil {
		t.Fatalf("assert.Nil: expected nil, got %v", got)
	}
}

func NotNil(t *testing.T, got any) {
	t.Helper()
	if got == nil {
		t.Fatalf("assert.NotNil: expected non-nil, got nil")
	}
}
