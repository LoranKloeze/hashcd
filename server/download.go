package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/cache"
	"github.com/lorankloeze/hashcd/config"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/middleware"
)

func Download(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := r.Context().Value(middleware.ContextRequestIdKey)
	ctx := log.WithLogger(r.Context(), log.L.WithField("reqid", id))

	hash, err := extractHash(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.G(ctx).Infof("Sending file %q", hash)

	path := filepath.Join(directoryTree(config.C.StorageDir, hash), hash)
	if !files.FileExists(path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reader, ok := cache.Retrieve(hash)
	if ok {
		log.G(ctx).Infof("Serving from cache %q", hash)
		w.Header().Set("X-Served-From", "cache on server")
		http.ServeContent(w, r, hash, time.Time{}, reader)
	} else {
		log.G(ctx).Infof("Serving from disk %q", hash)
		w.Header().Set("X-Served-From", "disk on server")
		cache.Insert(hash, path)
		http.ServeFile(w, r, path) // ServeFile sanitizes the path to prevent traversal attacks
	}
}

func extractHash(hashish string) (string, error) {
	re := regexp.MustCompile(`[a-f0-9]{64}`)
	res := re.FindString(hashish)
	if res == "" {
		return "", fmt.Errorf("could not extract hash")
	}
	return re.FindString(hashish), nil
}

func directoryTree(storageDir string, hash string) string {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := storageDir + "/"
	r := re.FindAllString(t, -1)
	return p + strings.Join(r, "/")
}
