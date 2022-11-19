package server

import (
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

	hash := p.ByName("hash")
	log.Infof("[%s] Sending file '%s'", id, hash)

	reader, ok := cache.GetFile(hash)
	if ok {
		log.Infof("[%s] Serving from cache", id)
		w.Header().Set("X-Served-From", "cache on server")
		http.ServeContent(w, r, hash, time.Time{}, reader)
	} else {
		log.Infof("[%s] Serving from disk", id)
		w.Header().Set("X-Served-From", "disk on server")
		path := filepath.Join(directoryTree(hash), hash)

		// ServeFile sanitizes the path to prevent traversal attacks
		http.ServeFile(w, r, path)
		cache.InsertFile(hash, path)
	}

	log.Infof("[%s] Sending finished", id)
}

func directoryTree(hash string) string {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := os.Getenv(envStorage) + "/"
	r := re.FindAllString(t, -1)
	return p + strings.Join(r, "/")
}
