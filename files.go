package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Errorf("Could not stat file '%s': %s", path, err)
		return 0, err
	}
	return fi.Size(), nil
}
