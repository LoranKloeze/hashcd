package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/middleware"
	log "github.com/sirupsen/logrus"
)

type fileStat struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

func HashList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(middleware.ContextRequestIdKey)
	res := []fileStat{}

	log.Infof("[%s] Sending list of files", id)

	w.Header().Set("Content-Type", "application/json")

	filepath.WalkDir(os.Getenv(envStorage), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Errorf("Could not walk directory: %s", err)
			os.Exit(1)
		}
		if !d.IsDir() {
			s, err := files.FileSize(path)
			if err != nil {
				s = 0
			}
			res = append(res, fileStat{Hash: d.Name(), Size: s})
		}
		return nil
	})

	json.NewEncoder(w).Encode(res)

	log.Infof("[%s] Sending finished", id)
}
