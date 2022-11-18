package cache

import (
	"bytes"
	"io"
	"os"

	"github.com/dgraph-io/ristretto"
	"github.com/lorankloeze/finalcd/files"
	log "github.com/sirupsen/logrus"
)

var cache *ristretto.Cache
var maxCacheSize = int64(1024 * 1024 * 512)   // 512 MB
var maxCacheItemSize = int64(1024 * 1024 * 3) // 3 MB

func Init() *ristretto.Cache {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,          // number of keys to track frequency of (10M).
		MaxCost:     maxCacheSize, // 10 MB
		BufferItems: 64,           // number of keys per Get buffer.
		OnEvict:     func(item *ristretto.Item) { log.Debugf("Cache: evicted %d - cost %d", item.Key, item.Cost) },
		Metrics:     false,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cache
}

func InsertFile(hash string, path string) {
	log.Debugf("Checking if '%s' needs to be cached", hash)

	s, err := files.FileSize(path)
	if err != nil {
		log.Errorf("Cannot determine file size, not caching '%s'", hash)
		return
	}

	if s > maxCacheItemSize {
		log.Debugf("File size (%d) exceeds cache item limit (%d), skipping cache", s, maxCacheItemSize)
		return
	}

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

	size, err := files.FileSize(path)
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

func GetFile(hash string) (io.ReadSeeker, bool) {
	log.Debugf("Searching cache entry for '%s'", hash)

	value, found := cache.Get(hash)
	if !found {
		log.Debugf("Cache entry not found for '%s'", hash)
		return nil, false
	}

	return bytes.NewReader(value.([]byte)), true
}
