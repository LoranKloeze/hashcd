package cache

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/ristretto"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/sizeutils"
	log "github.com/sirupsen/logrus"
)

var c *ristretto.Cache
var maxCacheItemSize int64

// Init initializes the cache and returns the cache. The cache size sets the cache
// capacity in megabytes. The item size sets the maximum size per item in megabytes.
func Init(cacheSize int64, itemSize int64) (*ristretto.Cache, error) {
	maxCacheItemSize = itemSize * sizeutils.Megabyte
	var err error

	c, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     cacheSize * sizeutils.Megabyte,
		BufferItems: 64,
		OnEvict:     func(item *ristretto.Item) { log.Debugf("Cache: evicted %d - cost %d", item.Key, item.Cost) },
		Metrics:     false,
	})
	if err != nil {
		return nil, fmt.Errorf("could not initialize cache: %v", err)
	}
	return c, nil
}

func Insert(hash string, path string) {
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

	contents, err := io.ReadAll(f) // ReadAll is fine here, it should fit in memory anyway
	if err != nil {
		log.Errorf("Could not read file for caching '%s': %s", hash, err)
		return
	}

	size, err := files.FileSize(path)
	if err != nil {
		log.Error("Cannot determine file size, not updating cache")
		return
	}
	ok := c.Set(hash, contents, size)
	if !ok {
		log.Error("Cache not updated")
	}

	log.Debugf("Finished creating cache entry for '%s'", hash)
}

func Retrieve(hash string) (io.ReadSeeker, bool) {
	log.Debugf("Searching cache entry for '%s'", hash)

	value, found := c.Get(hash)
	if !found {
		log.Debugf("Cache entry not found for '%s'", hash)
		return nil, false
	}

	return bytes.NewReader(value.([]byte)), true
}
