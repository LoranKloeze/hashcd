package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)

	store := InitStore("finalcd")
	defer store.Close()

	router := httprouter.New()
	router.POST("/", Upload)
	router.GET("/l", List)
	router.GET("/d/:hash", Download)

	n := negroni.New()
	n.Use(negroni.HandlerFunc(requestIdMiddleware))
	n.Use(negroni.HandlerFunc(corsMiddleWare))
	n.Use(negroni.HandlerFunc(logMiddleware))
	n.UseHandler(router)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}
