package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHashList(t *testing.T) {
	config.storageDir = tempStorageDir(t)
	defer os.RemoveAll(config.storageDir)

	createDummyHash(t, []byte{'a', 'b', 'c'})                // Hash:
	createDummyHash(t, []byte{'t', 'e', 's', 't'})           // Hash:
	createDummyHash(t, []byte{'g', 'o', 'g', 'o', 'g', 'o'}) // Hash:
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

func createDummyHash(t *testing.T, data []byte) string {
	hashData := sha256.New()
	io.Copy(hashData, bytes.NewBuffer(data))
	hashHex := hex.EncodeToString(hashData.Sum(nil))

	dir, err := initDirectories(hashHex, config.storageDir)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	f, err := os.Create(filepath.Join(dir, hashHex))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer f.Close()
	buf := bytes.NewBuffer(data)
	io.Copy(f, buf)

	return hashHex
}

func tempStorageDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "hashcd-storage")
	if err != nil {
		t.Fatal("cannot create temporary directory 'hashcd-storage'")
	}
	return tmpDir
}
