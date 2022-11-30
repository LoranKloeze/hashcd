package files

import (
	"os"

	"github.com/lorankloeze/hashcd/log"
)

func FileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		log.L.Errorf("Could not stat file %q: %s", path, err)
		return 0, err
	}
	return fi.Size(), nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
