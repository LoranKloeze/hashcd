package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/ristretto"
	log "github.com/sirupsen/logrus"
)

var cache *ristretto.Cache

func fileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Errorf("Could not stat file '%s': %s", path, err)
		return 0, err
	}
	return fi.Size(), nil
}

func InitStore(namespace string) *ristretto.Cache {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,                    // number of keys to track frequency of (10M).
		MaxCost:     1024 * 1024 * 1024 * 1, //  maximum cost of cache (1GB).
		BufferItems: 64,                     // number of keys per Get buffer.
		OnEvict:     func(item *ristretto.Item) { fmt.Printf("Evicted %d", item.Key) },
	})
	if err != nil {
		log.Fatal(err)
	}
	return cache
}

func InsertCacheFile(hash string, path string) {
	log.Debugf("Creating cache entry for '%s'", hash)

	f, err := os.Open(path)
	if err != nil {
		log.Errorf("Could not open file for caching '%s': %s", hash, err)
		return
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		log.Errorf("Could not read file for caching '%s': %s", hash, err)
		return
	}

	size, err := fileSize(path)
	if err != nil {
		log.Error("Cannot determine file size, not updating cache")
		return
	}
	ok := cache.Set(hash, contents, size)
	if !ok {
		log.Error("Cache not updated")
	}

	log.Debugf("Finished creating cache entry for '%s'", hash)
}

func FileFromCache(hash string) (io.ReadSeeker, bool) {
	log.Debugf("Searching cache entry for '%s'", hash)

	value, found := cache.Get(hash)
	if !found {
		log.Debugf("Cache entry not found for '%s'", hash)
		return nil, false
	}

	return bytes.NewReader(value.([]byte)), true
}
