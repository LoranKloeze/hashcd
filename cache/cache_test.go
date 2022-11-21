package cache

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestInit(t *testing.T) {
	c, err := Init(10, 0)

	if err != nil {
		t.Fatalf("unexpected error setting up cache: %v", err)
	}
	defer c.Close()

	want := int64(10485760)
	if c.MaxCost() != want {
		t.Errorf("expected MaxCost() = %v, got %v", want, c.MaxCost())
	}
}

func TestInsert(t *testing.T) {
	c, err := Init(10, 2)
	if err != nil {
		t.Fatalf("unexpected error setting up cache: %v", err)
	}
	defer c.Close()

	f := dummyFile(t, 10)
	defer cleanDummyFile(f)

	hash := "myhash"
	Insert(hash, f.Name())

	time.Sleep(10 * time.Millisecond) // Wait for value to pass through buffers

	val, found := c.Get(hash)
	if !found {
		t.Errorf("Get() didn't find hash %q", hash)
	}

	f.Seek(0, 0)
	want, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("unexpected error reading file: %v", err)
	}
	if !cmp.Equal(val.([]byte), want) {
		t.Errorf("expected cache content to equal file contents: %v != %v", val.([]byte), want)
	}

}

func TestRetrieve(t *testing.T) {
	c, err := Init(10, 2)
	if err != nil {
		t.Fatalf("unexpected error setting up cache: %v", err)
	}
	defer c.Close()

	f := dummyFile(t, 10)
	defer cleanDummyFile(f)

	fContent, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("unexpected error reading file: %v", err)
	}

	hash := "myhash"
	ok := c.Set(hash, fContent, 10)
	if !ok {
		t.Fatalf("unexpected error setting cache: %v", err)
	}

	time.Sleep(10 * time.Millisecond) // Wait for value to pass through buffers

	r, ok := Retrieve(hash)
	if !ok {
		t.Fatalf("expected GetFile to find cache content for %q", hash)
	}
	cContent, err := io.ReadAll(r)
	if !ok {
		t.Fatalf("unexpected error reading cache: %v", err)
	}
	if !cmp.Equal(cContent, fContent) {
		t.Errorf("expected cache content to equal file contents: %v != %v", cContent, fContent)
	}

}

func TestInsertTooBig(t *testing.T) {
	c, err := Init(10, 1)
	if err != nil {
		t.Fatalf("unexpected error setting up cache: %v", err)
	}
	defer c.Close()

	f := dummyFile(t, 1024*1024*1+1)
	defer cleanDummyFile(f)

	hash := "myhash"
	Insert(hash, f.Name())

	time.Sleep(10 * time.Millisecond) // Wait for value to pass through buffers

	_, found := c.Get(hash)
	if found {
		t.Errorf("Get() should not return a value for %q", hash)
	}

}

func dummyFile(t *testing.T, size int) *os.File {
	f, err := os.CreateTemp("/tmp", "hashcd")
	if err != nil {
		t.Fatalf("Setup failed: could not create temp file: %v", err)
	}

	f.Write(bytes.Repeat([]byte{'x'}, int(size)))

	if err != nil {
		t.Fatalf("Setup failed: could not write temp file: %v", err)
	}

	return f
}

func cleanDummyFile(f *os.File) {
	f.Close()
	os.Remove(f.Name())
}
