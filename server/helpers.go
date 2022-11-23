package server

import (
	"os"
	"path/filepath"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// initDirectories creates a directory structure for a given hash in
// the storage directory
//
// Example: a36b7e89dd creates <storageDir>/a3/6b/7e/89 (the last 2 characters are skipped)
func initDirectories(hash string, storageDir string) (createdPath string, err error) {
	if storageDir == "" {
		log.Fatalf("storageDir cannot be empty")
	}

	t := hash[0 : len(hash)-2] // We don't create the last subdirectory

	re := regexp.MustCompile(`..`)
	r := []string{storageDir}
	r = append(r, re.FindAllString(t, -1)...)
	p := filepath.Join(r...)

	err = os.MkdirAll(p, 0755)
	if err != nil {
		log.Errorf("Could not create directory storage tree: %s", err)
		return "", err
	}

	log.Debugf("Initialized directory storage tree '%s'", p)

	return p, nil
}

func validateConfig() {
	if config.storageDir == "" {
		log.Fatal("config.storageDir cannot be empty")
	}
}
