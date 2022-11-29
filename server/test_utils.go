package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func createDummyHash(t *testing.T, data []byte) string {
	hashData := sha256.New()
	io.Copy(hashData, bytes.NewBuffer(data))
	hashHex := hex.EncodeToString(hashData.Sum(nil))

	dir, err := initDirectories(hashHex, Config.StorageDir)
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
