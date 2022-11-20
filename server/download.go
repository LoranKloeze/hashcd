package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/finalcd/cache"
	"github.com/lorankloeze/finalcd/middleware"
	log "github.com/sirupsen/logrus"
)

func Download(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(middleware.ContextRequestIdKey)

	hash, err := extractHash(p.ByName("hashish"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Printf("La: %q\n", hash)
	log.Infof("[%s] Sending file '%s'", id, hash)

	path := filepath.Join(directoryTree(hash), hash)
	if !fileExists(path) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reader, ok := cache.GetFile(hash)
	if ok {
		log.Infof("[%s] Serving from cache", id)
		w.Header().Set("X-Served-From", "cache on server")
		http.ServeContent(w, r, hash, time.Time{}, reader)
	} else {
		log.Infof("[%s] Serving from disk", id)
		w.Header().Set("X-Served-From", "disk on server")
	}

	http.ServeFile(w, r, path) // ServeFile sanitizes the path to prevent traversal attacks
	cache.InsertFile(hash, path)

	log.Infof("[%s] Sending finished", id)
}

func extractHash(hashish string) (string, error) {
	re := regexp.MustCompile(`[a-f0-9]{64}`)
	res := re.FindString(hashish)
	if res == "" {
		return "", fmt.Errorf("could not extract hash")
	}
	return re.FindString(hashish), nil
}

func directoryTree(hash string) string {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := os.Getenv(envStorage) + "/"
	r := re.FindAllString(t, -1)
	return p + strings.Join(r, "/")
}
