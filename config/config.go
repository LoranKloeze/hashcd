package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lorankloeze/hashcd/files"
	"github.com/sirupsen/logrus"
)

// C contains the current loaded configuration
var C Configuration

// Configuration describes the data model of the application configuration
type Configuration struct {

	// CacheSize is the maxium memory in MiB used by the storage cache
	CacheSize int64

	// CacheItemSize is the size cutoff in MiB where files are not cached anymore
	CacheItemSize int64

	// ListenAddr is the address the server listens on e.g. 127.0.0.1:8080
	ListenAddr string

	// LogLevel is the log level used by the application logger
	LogLevel logrus.Level

	// StorageDir is the directory where uploaded files are stored
	StorageDir string
}

const (
	envCacheSize     = "HASHCD_CACHE_SIZE"
	envCacheItemSize = "HASHCD_CACHE_ITEM_SIZE"
	envStorageDir    = "HASHCD_STORAGE_DIR"
	envListenAddr    = "HASHCD_LISTEN_ADDR"
	envLogLevel      = "HASHCD_LOGLEVEL"
)

var defaultCfg Configuration = Configuration{
	CacheSize:     128,
	CacheItemSize: 2,
	LogLevel:      logrus.InfoLevel,
	ListenAddr:    "127.0.0.1:8080",
	StorageDir:    "", // not valid but checked again in Load()
}

// Load initializes the configuration from the environment
func Load() error {
	C = defaultCfg
	var err error

	// CacheSize
	if v, ok := os.LookupEnv(envCacheSize); ok {
		C.CacheSize, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("environment variable %q = %q is not a number", envCacheSize, v)
		}
	}

	// CacheItemSize
	if v, ok := os.LookupEnv(envCacheItemSize); ok {
		C.CacheItemSize, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("environment variable %q = %q is not a number", envCacheItemSize, v)
		}
	}

	// StorageDir
	if v, ok := os.LookupEnv(envStorageDir); ok {
		if files.FileExists(v) {
			C.StorageDir = v
		} else {
			return fmt.Errorf("environment variable %q = %q is not a valid or existing directory", envStorageDir, v)
		}
	} else {
		return fmt.Errorf("environment variable %q is not set", envStorageDir)
	}

	// ListenAddr
	if v, ok := os.LookupEnv(envListenAddr); ok {
		C.ListenAddr = v
	}

	// LogLevel
	lala, ok := os.LookupEnv(envLogLevel)
	fmt.Println(lala, ok)
	if v, ok := os.LookupEnv(envLogLevel); ok {
		lvl, err := logrus.ParseLevel(v)
		if err != nil {
			return fmt.Errorf("environment variable %q = %q is not a valid log level", envLogLevel, v)
		}
		C.LogLevel = lvl
	}

	return nil
}
