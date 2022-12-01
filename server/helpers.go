package server

import (
	"context"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lorankloeze/hashcd/log"
)

// initDirectories creates a directory structure for a given hash in
// the storage directory
//
// Example: a36b7e89dd creates <storageDir>/a3/6b/7e/89 (the last 2 characters are skipped)
func initDirectories(hash string, storageDir string) (createdPath string, err error) {
	ctx := log.WithLogger(context.Background(), log.L.WithField("hash", hash))

	if storageDir == "" {
		log.G(ctx).Fatalf("storageDir cannot be empty in initDirectories")
	}

	t := hash[0 : len(hash)-2] // We don't create the last subdirectory

	re := regexp.MustCompile(`..`)
	r := []string{storageDir}
	r = append(r, re.FindAllString(t, -1)...)
	p := filepath.Join(r...)

	err = os.MkdirAll(p, 0755)
	if err != nil {
		log.G(ctx).Errorf("Failed to create directory storage tree: %v", err)
		return "", err
	}
	log.G(ctx).Debugf("Initialized directory storage tree %q", p)

	return p, nil
}
