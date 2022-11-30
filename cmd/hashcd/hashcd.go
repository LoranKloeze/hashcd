package main

import (
	"net/http"
	"os"
	"strconv"

	_ "net/http/pprof"

	"github.com/dgraph-io/ristretto"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/cache"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/middleware"
	"github.com/lorankloeze/hashcd/server"
	"github.com/sirupsen/logrus"

	"github.com/urfave/negroni"
)

const envCacheSize = "HASHCD_CACHE_SIZE"
const envCacheItemSize = "HASHCD_CACHE_ITEM_SIZE"

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		log.L.Info("No .env file found, that's fine: using regular environment")
	}

	server.Config = server.Configuration{
		StorageDir: os.Getenv("HASHCD_STORAGE"),
	}

	c := initCache()
	defer c.Close()

	router := httprouter.New()
	router.POST("/u", server.Upload)
	router.GET("/l", server.HashList)
	router.GET("/d/:hashish", server.Download)

	n := negroni.New()
	n.Use(negroni.HandlerFunc(middleware.RequestId))
	n.Use(negroni.HandlerFunc(middleware.Cors))
	n.Use(negroni.HandlerFunc(middleware.Log))
	n.UseHandler(router)

	log.L.Infof("PID: %d\n", os.Getpid())
	log.L.Info("Listening on port 8080")
	log.L.Fatal(http.ListenAndServe(":8080", n))
}

func initCache() *ristretto.Cache {

	cacheSize, err := strconv.ParseInt(os.Getenv(envCacheSize), 10, 64)
	if err != nil {
		log.L.Fatalf("envCacheSize: %q is not a number", envCacheSize)
	}

	maxCacheItemSize, err := strconv.ParseInt(os.Getenv(envCacheItemSize), 10, 64)
	if err != nil {
		log.L.Fatalf("envCacheItemSize: %q is not a number", envCacheItemSize)
	}

	c, err := cache.Init(cacheSize, maxCacheItemSize)
	if err != nil {
		log.L.Fatalf("Failed to setup cache: %v", err)
	}

	return c
}
