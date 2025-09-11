package decorator

import "testing"

func TestGet(t *testing.T) {
	elems, err := GetElements(&Options{Context: "privatperson"})
	if err != nil {
		t.Fatal(err)
	}
	if elems.HeadAssets == "" {
		t.Error("HeadAssets is empty")
	}
	if elems.Header == "" {
		t.Error("Header is empty")
	}
	if elems.Footer == "" {
		t.Error("Footer is empty")
	}
	if elems.Scripts == "" {
		t.Error("Scripts is empty")
	}
}
