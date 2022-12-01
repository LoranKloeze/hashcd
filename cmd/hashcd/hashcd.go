package main

import (
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/dgraph-io/ristretto"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/hashcd/cache"
	"github.com/lorankloeze/hashcd/config"
	"github.com/lorankloeze/hashcd/log"
	"github.com/lorankloeze/hashcd/middleware"
	"github.com/lorankloeze/hashcd/server"
	"github.com/sirupsen/logrus"

	"github.com/urfave/negroni"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	err := godotenv.Load()
	if err != nil {
		log.L.Info("No .env file found, that's fine: using regular environment")
	}

	err = config.Load()
	if err != nil {
		log.L.Fatalf("Could not load configuration: %v", err)
	}

	logrus.SetLevel(config.C.LogLevel)

	c := initCache(config.C.CacheSize, config.C.CacheItemSize)
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

func initCache(cacheSize, maxCacheItemSize int64) *ristretto.Cache {

	c, err := cache.Init(cacheSize, maxCacheItemSize)
	if err != nil {
		log.L.Fatalf("Failed to setup cache: %v", err)
	}

	return c
}
