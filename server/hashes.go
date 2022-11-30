package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/middleware"
)

type fileStat struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

func HashList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	validateConfig()
	id := r.Context().Value(middleware.ContextRequestIdKey)
	ctx := log.WithLogger(r.Context(), log.L.WithField("reqid", id))

	res := []fileStat{}

	log.G(ctx).Info("Sending file list")

	w.Header().Set("Content-Type", "application/json")

	filepath.WalkDir(Config.StorageDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.G(ctx).Fatalf("Failed to walk directory: %v", err)
		}

		if !d.IsDir() {
			s, err := files.FileSize(path)
			if err != nil {
				log.G(ctx).Fatalf("Failed to determine file size: %v", err)
			}
			res = append(res, fileStat{Hash: d.Name(), Size: s})
		}
		return nil
	})

	json.NewEncoder(w).Encode(res)
}
