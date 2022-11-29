package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/middleware"
	"github.com/lorankloeze/hashcd/sizeutils"
	log "github.com/sirupsen/logrus"
)

type okResponse struct {
	Hash string `json:"hash"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	validateConfig()
	id := r.Context().Value(middleware.ContextRequestIdKey)

	log.Infof("[%s] Receiving file", id)

	err := r.ParseMultipartForm(10 * sizeutils.Megabyte)
	if err != nil {
		log.Printf("[%s] Could not parse form: %s\n", id, err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse{"Could not parse form data"})
		return
	}

	field, ok := r.MultipartForm.File["f"]
	if !ok {
		log.Errorf("[%s] Form field 'f' not found", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse{"Form field 'f' not found"})
		return
	}
	f, err := field[0].Open()
	if err != nil {
		log.Errorf("[%s] Could not open file from form: %s\n", id, err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse{"Could not open file from form"})
		return
	}
	defer f.Close()

	hash := genHash(f)
	dirs, err := initDirectories(hash, Config.StorageDir)
	if err != nil {
		log.Fatalf("[%s] Could not create directory storage tree %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("%s/%s", dirs, hash)
	if fileExists(path) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(okResponse{hash})
		return
	}

	f.Seek(0, 0)

	f1, err := os.Create(path)
	if err != nil {
		log.Fatalf("[%s] Could not create file: %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f1.Close()

	written, err := io.Copy(f1, f)
	if err != nil {
		log.Fatalf("[%s] Could not write to file: %s\n", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debugf("[%s] Written %d bytes to %s", id, written, path)
	log.Infof("[%s] Saved file %s", id, hash)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res := okResponse{hash}
	json.NewEncoder(w).Encode(res)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func genHash(f io.Reader) string {
	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		log.Fatalf("Could not calculate SHA256 hash %s\n", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
