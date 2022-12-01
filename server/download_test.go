package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lorankloeze/hashcd/cache"
	"github.com/lorankloeze/hashcd/config"
)

func TestDownload(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)

	exp := []byte{'g', 'o', 'g', 'o', 'g', 'o'}
	h := createDummyHash(t, exp)
	p := fmt.Sprintf("/d/%s", h)
	r := httptest.NewRequest("GET", p, nil)
	defer r.Body.Close()
	w := httptest.NewRecorder()

	Download(w, r, nil)

	got, _ := io.ReadAll(w.Body)
	if !cmp.Equal(exp, got) {
		t.Errorf("expected download to return %v, got %v", exp, got)
	}

}

func TestDownloadNonExisting(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)

	r := httptest.NewRequest("GET", "/d/idonotexist", nil)
	defer r.Body.Close()
	w := httptest.NewRecorder()

	Download(w, r, nil)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected non existing file to return 404")
	}
}

func TestDownloadCachedFile(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)
	c, _ := cache.Init(1, 1)
	defer c.Close()

	h := createDummyHash(t, []byte{'g', 'o', 'g', 'o', 'g', 'o'})
	p := fmt.Sprintf("/d/%s", h)

	// First download serves from disk
	r := httptest.NewRequest("GET", p, nil)
	defer r.Body.Close()
	w := httptest.NewRecorder()
	Download(w, r, nil)

	exp := "disk on server"
	if w.Result().Header["X-Served-From"][0] != exp {
		t.Errorf("first download, expected a file from %q, got a file from %q", exp, w.Result().Header["X-Served-From"][0])
	}

	time.Sleep(10 * time.Millisecond) // Wait for value to pass through cache buffers

	// Second download serves from cache
	r = httptest.NewRequest("GET", p, nil)
	w = httptest.NewRecorder()
	Download(w, r, nil)

	exp = "cache on server"
	if w.Result().Header["X-Served-From"][0] != exp {
		t.Errorf("second download, expected a file from %q, got a file from %q", exp, w.Result().Header["X-Served-From"][0])
	}

}

func TestDownloadUncachedFile(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)
	c, _ := cache.Init(1, 1)
	defer c.Close()

	mb := 1024 * 1024
	b := make([]byte, 2*mb)
	h := createDummyHash(t, b)
	p := fmt.Sprintf("/d/%s", h)

	// First download serves from disk
	r := httptest.NewRequest("GET", p, nil)
	defer r.Body.Close()
	w := httptest.NewRecorder()
	Download(w, r, nil)

	exp := "disk on server"
	if w.Result().Header["X-Served-From"][0] != exp {
		t.Errorf("first download, expected a file from %q, got a file from %q", exp, w.Result().Header["X-Served-From"][0])
	}

	time.Sleep(10 * time.Millisecond) // Wait for value to pass through cache buffers

	// Second download serves from disk too
	r = httptest.NewRequest("GET", p, nil)
	w = httptest.NewRecorder()
	Download(w, r, nil)

	exp = "disk on server"
	if w.Result().Header["X-Served-From"][0] != exp {
		t.Errorf("second download, expected a file from %q, got a file from %q", exp, w.Result().Header["X-Served-From"][0])
	}
}
