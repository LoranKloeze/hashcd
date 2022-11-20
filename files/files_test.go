package files

import (
	"os"
	"testing"
)

func TestExistingFileSize(t *testing.T) {
	want := int64(25)
	f := dummyFile(t, want)
	defer cleanDummyFile(f)

	s, err := FileSize(f.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if s != want {
		t.Errorf("FileSize(%q) == %v, want %v", f.Name(), s, want)
	}

}

func TestNonExistingFileSize(t *testing.T) {
	_, err := FileSize("/i/most/certainly/do/not/exist/4256")
	if err == nil {
		t.Errorf("FileSize() with non existing file should return an error")
	}
}

func dummyFile(t *testing.T, size int64) *os.File {
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
