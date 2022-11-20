package main

import (
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/lorankloeze/finalcd/cache"
	"github.com/lorankloeze/finalcd/middleware"
	"github.com/lorankloeze/finalcd/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info("No .env file found, that's fine: using regular environment")
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)

	c := cache.Init()
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
