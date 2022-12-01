package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/log"
	"github.com/sirupsen/logrus"
)

// C contains the current loaded configuration
var C Configuration

// Configuration describes the data model of the application configuration
type Configuration struct {

	// CacheSize is the maxium memory in megabytes used by the storage cache
	CacheSize int64

	// CacheItemSize is the size cutoff where files are not cached anymore
	CacheItemSize int64

	// LogLevel is the log level used by the application logger
	LogLevel logrus.Level

	// StorageDir is the directory where uploaded files are stored
	StorageDir string
}

const (
	envCacheSize     = "HASHCD_CACHE_SIZE"
	envCacheItemSize = "HASHCD_CACHE_ITEM_SIZE"
	envStorageDir    = "HASHCD_STORAGE_DIR"
)

var defaultCfg Configuration = Configuration{
	CacheSize:     128,
	CacheItemSize: 2,
	LogLevel:      logrus.InfoLevel,
	StorageDir:    "",
}

// Load initializes the configuration from the environment
func Load() error {
	C = defaultCfg

	cacheSize, err := strconv.ParseInt(os.Getenv(envCacheSize), 10, 64)
	if err != nil {
		log.L.Warnf("Cache size: %q is not a number, using default", os.Getenv(envCacheSize))
	} else {
		C.CacheSize = cacheSize
	}

	cacheItemSize, err := strconv.ParseInt(os.Getenv(envCacheItemSize), 10, 64)
	if err != nil {
		log.L.Warnf("Cache item size: %q is not a number, using default", os.Getenv(envCacheItemSize))
	} else {
		C.CacheItemSize = cacheItemSize
	}

	if !files.FileExists(os.Getenv(envStorageDir)) {
		return fmt.Errorf("storage directory: %q does not exist", os.Getenv(envStorageDir))
	}
	C.StorageDir = os.Getenv(envStorageDir)

	return nil

}
