package assert

import "testing"

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

func NotEqual[T comparable](t *testing.T, got, unexpected T) {
	t.Helper()
	if got == unexpected {
		t.Errorf("assert.NotEqual: got unexpected value %v", unexpected)
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
