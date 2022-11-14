package main

import (
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func main() {
	// defer profile.Start(profile.MemProfile).Stop()
	// go func(c chan int) {
	// 	http.ListenAndServe(":8081", nil)
	// }()

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

	log.Printf("PID: %d\n", os.Getpid())
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", n))
}
