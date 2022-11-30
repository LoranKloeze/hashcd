package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/middleware"

	"github.com/lorankloeze/hashcd/sizeutils"
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
	ctx := log.WithLogger(r.Context(), log.L.WithField("reqid", id))

	log.G(ctx).Info("Receiving file")

	file, err := fileFromForm(r)
	if err != nil {
		log.G(ctx).Errorf("Unprocessable form: %v", err)
		respondInvalid(w, err)
		return
	}
	defer file.Close()

	hash := genHash(ctx, file)
	dirs, err := initDirectories(hash, Config.StorageDir)
	if err != nil {
		log.G(ctx).Errorf("Directory storage tree not created: %v", err)
		respondError(w)
		return
	}

	path := filepath.Join(dirs, hash)
	if files.FileExists(path) {
		respondHash(w, hash, http.StatusOK)
		return
	}

	err = saveFile(ctx, file, hash, path)
	if err != nil {
		log.G(ctx).Errorf("File not saved: %v", err)
		respondError(w)
		return
	}
	respondHash(w, hash, http.StatusCreated)
}

func genHash(ctx context.Context, f io.Reader) string {
	hash := sha256.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		log.G(ctx).Fatalf("Failed to calculate SHA256 hash: %v", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func fileFromForm(r *http.Request) (multipart.File, error) {
	err := r.ParseMultipartForm(10 * sizeutils.Megabyte)
	if err != nil {
		return nil, fmt.Errorf("failed to parse form data")
	}

	field, ok := r.MultipartForm.File["f"]
	if !ok {
		return nil, fmt.Errorf("form field 'f' not found")
	}

	f, err := field[0].Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file from form")
	}

	return f, nil
}

func saveFile(ctx context.Context, srcFile multipart.File, hash string, path string) error {
	srcFile.Seek(0, 0)

	dstFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	log.G(ctx).Infof("Saved file %q to disk", hash)
	return nil
}

func respondHash(w http.ResponseWriter, hash string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(okResponse{hash})
}

func respondError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func respondInvalid(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(errorResponse{fmt.Sprintf("%v", err)})
}
