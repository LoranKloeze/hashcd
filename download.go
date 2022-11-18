package main

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func directoryTree(hash string) string {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := "/home/loran/git/lab/finalcd/storage/"
	r := re.FindAllString(t, -1)
	return p + strings.Join(r, "/")
}

func Download(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(contextRequestIdKey)

	hash := p.ByName("hash")
	log.Infof("[%s] Sending file '%s'", id, hash)

	reader, ok := FileFromCache(hash)
	if ok {
		log.Infof("[%s] Serving from cache", id)
		w.Header().Set("Served-From", "cache on server")
		http.ServeContent(w, r, hash, time.Time{}, reader)
	} else {
		log.Infof("[%s] Serving from disk", id)
		w.Header().Set("Served-From", "disk on server")
		path := filepath.Join(directoryTree(hash), hash)

		// ServeFile sanitizes the path to prevent traversal attacks
		http.ServeFile(w, r, path)
		InsertCacheFile(hash, path)
	}

	log.Infof("[%s] Sending finished", id)
}
