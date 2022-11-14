package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type FileStat struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

func List(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(contextRequestIdKey)
	var res []FileStat

	log.Infof("[%s] Sending list of files", id)

	w.Header().Set("Content-Type", "application/json")

	filepath.WalkDir("/home/loran/git/lab/finalcd/storage/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			s, err := FileSize(path)
			if err != nil {
				s = 0
			}
			res = append(res, FileStat{Hash: d.Name(), Size: s})
		}
		return nil
	})

	json.NewEncoder(w).Encode(res)

	log.Infof("[%s] Sending finished", id)
}
