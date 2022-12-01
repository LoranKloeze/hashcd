package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lorankloeze/hashcd/config"
)

func TestUpload(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)

	w := httptest.NewRecorder()
	r, expHash := makeUploadRequest(t, []byte{'a', 'b', 'c'})
	defer r.Body.Close()

	Upload(w, r, nil)

	// If hash is abcdef1011, directory is storage/ab/cd/ef/10 (thus without the last 2 characters)
	re := regexp.MustCompile(`..`)
	d := re.FindAllString(expHash, -1)
	dir := filepath.Join(config.C.StorageDir, filepath.Join(d[0:len(d)-1]...), expHash)

	f, err := os.Open(dir)
	if err != nil {
		t.Errorf("cannot find uploaded file at %q", dir)
	}
	defer f.Close()

	got, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	want := []byte{'a', 'b', 'c'}
	if !cmp.Equal(got, want) {
		t.Errorf("contents of disk file != uploaded contents: got %v, want %v", got, want)
	}

}

func makeUploadRequest(t *testing.T, data []byte) (req *http.Request, hash string) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("f", "xyz.xyz")
		if err != nil {
			t.Error(err)
		}
		part.Write(data)
	}()
	r := httptest.NewRequest("POST", "/", pr)

	hashData := sha256.New()
	io.Copy(hashData, bytes.NewBuffer(data))
	hashHex := hex.EncodeToString(hashData.Sum(nil))

	r.Header.Add("Content-Type", writer.FormDataContentType())
	return r, hashHex
}
