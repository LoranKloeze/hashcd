package server

import (
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDownload(t *testing.T) {
	Config.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(Config.StorageDir)

	exp := []byte{'g', 'o', 'g', 'o', 'g', 'o'}
	h := createDummyHash(t, exp)
	p := fmt.Sprintf("/d/%s", h)
	r := httptest.NewRequest("GET", p, nil)
	w := httptest.NewRecorder()

	Download(w, r, nil)

	got, _ := io.ReadAll(w.Body)
	if !cmp.Equal(exp, got) {
		t.Errorf("expected download to return %v, got %v", exp, got)
	}

}
