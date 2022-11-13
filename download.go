package main

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func directoryTree(hash string) string {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := "/home/loran/git/lab/mycdn/storage/"
	r := re.FindAllString(t, -1)
	return p + strings.Join(r, "/")
}

func Download(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(contextRequestIdKey)

	hash := p.ByName("hash")
	log.Infof("[%s] Sending file '%s'", id, hash)

	path := filepath.Join(directoryTree(hash), hash)

	// ServeFile sanitizes the path to prevent traversal attacks
	http.ServeFile(w, r, path)

	log.Infof("[%s] Sending finished", id)
}
