package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type okResponse struct {
	Hash string `json:"hash"`
}

type errorResponse struct {
	Error string `json:"error"`
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

func initDirectories(hash string) (string, error) {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := "/home/loran/git/lab/mycdn/storage/"
	r := re.FindAllString(t, -1)
	p += strings.Join(r, "/")
	err := os.MkdirAll(p, 0755)
	if err != nil {
		log.Errorf("Could not create directory storage tree: %s", err)
		return "", err
	}
	log.Debugf("Initialized directory storage tree '%s'", p)
	return p, nil
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(contextRequestIdKey)
	maxMem := int64(10 << 20) // 10MB

	log.Infof("[%s] Receiving file", id)
	err := r.ParseMultipartForm(maxMem)
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
	dirs, err := initDirectories(hash)
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
