package cache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/ristretto"
	"github.com/lorankloeze/hashcd/files"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/sizeutils"
)

const logHashKey = "hash"

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
		OnEvict:     onEvict,
		Metrics:     false,
	})
	if err != nil {
		return nil, fmt.Errorf("could not initialize cache: %v", err)
	}
	return c, nil
}

func Insert(hash string, path string) {
	ctx := log.WithLogger(context.Background(), log.L.WithField(logHashKey, hash))
	log.G(ctx).Debugf("Check if caching is needed")

	s, err := files.FileSize(path)
	if err != nil {
		log.G(ctx).Errorf("Failed to determine file size, skipping cache: %v", err)
		return
	}

	if s > maxCacheItemSize {
		log.G(ctx).Debugf("File size (%d) exceeds cache item limit (%d), skipping cache", s, maxCacheItemSize)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		log.G(ctx).Errorf("Could not open file for caching: %v", hash, err)
		return
	}
	defer f.Close()

	contents, err := io.ReadAll(f) // ReadAll is fine here, it should fit in memory anyway
	if err != nil {
		log.G(ctx).Errorf("Could not read file for caching: %v", err)
		return
	}

	size, err := files.FileSize(path)
	if err != nil {
		log.G(ctx).Errorf("Cannot determine file size, not updating cache: %v", err)
		return
	}
	ok := c.Set(hash, contents, size)
	if !ok {
		log.G(ctx).Errorf("Cache not updated: Ristretto Set method returned false")
		return
	}

	log.G(ctx).Debugf("Finished creating cache entry")
}

func Retrieve(hash string) (io.ReadSeeker, bool) {
	ctx := log.WithLogger(context.Background(), log.L.WithField(logHashKey, hash))
	log.G(ctx).Debugf("Searching cache entry")

	value, found := c.Get(hash)
	if !found {
		log.G(ctx).Debugf("Cache entry not found")
		return nil, false
	}

	return bytes.NewReader(value.([]byte)), true
}

func onEvict(item *ristretto.Item) {
	ctx := log.WithLogger(context.Background(), log.L.WithField(logHashKey, item.Key))
	log.G(ctx).Debugf("Evicted from cache - cost %d", item.Cost)
}
