package server

import (
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUpload(t *testing.T) {
	Config.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(Config.StorageDir)

	// Setup form
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("f", "afile.txt")
		if err != nil {
			t.Error(err)
		}
		part.Write([]byte{'a', 'b', 'c'})
	}()

	// Setup request
	r := httptest.NewRequest("POST", "/", pr)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	Upload(w, r, nil)

	expHash := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	re := regexp.MustCompile(`..`)

	d := re.FindAllString(expHash, -1)
	dir := filepath.Join(Config.StorageDir, filepath.Join(d[0:len(d)-1]...), expHash)
	f, err := os.Open(dir)
	if err != nil {
		t.Errorf("cannot find uploaded file at %q", dir)
	}
	defer f.Close()
	contents, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	want := []byte{'a', 'b', 'c'}
	if !cmp.Equal(contents, want) {
		t.Errorf("contents of disk file != uploaded contents: %v != %v", contents, want)
	}

}
