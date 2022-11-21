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
	"github.com/lorankloeze/hashcd/middleware"
	"github.com/lorankloeze/hashcd/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

const envCacheSize = "HASHCD_CACHE_SIZE"
const envCacheItemSize = "HASHCD_CACHE_ITEM_SIZE"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found, that's fine: using regular environment")
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)

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

	log.Printf("PID: %d\n", os.Getpid())
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}

func initCache() *ristretto.Cache {

	cacheSize, err := strconv.ParseInt(os.Getenv(envCacheSize), 10, 64)
	if err != nil {
		log.Errorf("%s is not a number: %s", envCacheSize, err)
	}

	maxCacheItemSize, err := strconv.ParseInt(os.Getenv(envCacheItemSize), 10, 64)
	if err != nil {
		log.Errorf("%s is not a number: %s", envCacheItemSize, err)
	}

	c, err := cache.Init(cacheSize, maxCacheItemSize)
	if err != nil {
		log.Fatalf("Error setting up cache: %v", err)
	}

	return c
}
