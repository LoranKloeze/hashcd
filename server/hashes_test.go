package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lorankloeze/hashcd/config"
)

func TestHashList(t *testing.T) {
	config.C.StorageDir = tempStorageDir(t)
	defer os.RemoveAll(config.C.StorageDir)

	createDummyHash(t, []byte{'a', 'b', 'c'})
	createDummyHash(t, []byte{'t', 'e', 's', 't'})
	createDummyHash(t, []byte{'g', 'o', 'g', 'o', 'g', 'o'})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	HashList(w, r, nil)
	defer r.Body.Close()

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Result().StatusCode)
	}

	var res []fileStat
	json.NewDecoder(w.Body).Decode(&res)

	if len(res) != 3 {
		t.Errorf("expected resonse json to have 3 items, got %d", len(res))
		return
	}

	for i, e := range []fileStat{
		{Hash: "63cb946beb677e9f8d28ec5a8c6f6a929eb5b7dddc4b286f86345813c2d58e5a", Size: 6},
		{Hash: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", Size: 4},
		{Hash: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad", Size: 3},
	} {
		if res[i].Hash != e.Hash {
			t.Errorf("expcted hash %q for item %d, got %q", e.Hash, i, res[i].Hash)
		}
		if res[i].Size != e.Size {
			t.Errorf("expcted size %d for item %d, got %d", e.Size, i, res[i].Size)
		}
	}
}
