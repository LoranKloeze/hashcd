package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var redisCli *redis.Client
var redisNamespace string

func InitStore(namespace string) *redis.Client {
	redisNamespace = namespace
	redisCli = redis.NewClient(&redis.Options{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Addr:         "localhost:8910",
		Password:     "",
		DB:           0,
	})
	return redisCli
}

func InsertCacheFile(hash string, path string) {
	log.Debugf("Creating cache entry for '%s'", hash)

	key := fmt.Sprintf("%s:f:%s", redisNamespace, hash)
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("Could not open file for caching '%s': %s", hash, err)
		return
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		log.Errorf("Could not read file for caching '%s': %s", hash, err)
		return
	}

	err = redisCli.Set(context.Background(), key, contents, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Finished creating cache entry for '%s'", hash)
}

func FileFromCache(hash string) (io.ReadSeeker, bool) {
	log.Debugf("Searching cache entry for '%s'", hash)

	key := fmt.Sprintf("%s:f:%s", redisNamespace, hash)
	c, err := redisCli.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, false
	}
	return bytes.NewReader(c), true
}
