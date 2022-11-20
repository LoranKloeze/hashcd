package files

import (
	"os"
	"testing"
)

func TestFileSize(t *testing.T) {
	want := int64(25)
	f := createDummyFile(t, want)
	defer cleanDummyFile(f)

	s, err := FileSize(f.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if s != want {
		t.Errorf("FileSize(%q) == %v, want %v", f.Name(), s, want)
	}

}

func createDummyFile(t *testing.T, size int64) *os.File {
	f, err := os.CreateTemp("/tmp", "hashcd")
	if err != nil {
		t.Fatalf("Setup failed: could not create temp file: %v", err)
	}

	s := ""
	for i := int64(0); i < size; i++ {
		s += "x"
	}

	_, err = f.WriteString(s)
	if err != nil {
		t.Fatalf("Setup failed: could not write temp file: %v", err)
	}

	return f
}

func cleanDummyFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}
