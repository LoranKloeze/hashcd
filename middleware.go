package main

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

type ContextKey string

const contextRequestIdKey ContextKey = "requestId"

func clientIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

func requestIdMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := context.WithValue(r.Context(), contextRequestIdKey, uuid.New())
	r = r.WithContext(ctx)
	next(rw, r)
}

func logMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tStart := time.Now()
	id := r.Context().Value(contextRequestIdKey)

	log.Infof("[%s] Request start| %s | %s | %s", id, clientIp(r), r.Method, r.URL.Path)
	lrw := negroni.NewResponseWriter(rw)
	next(lrw, r)
	dur := time.Since(tStart)
	log.Infof("[%s] Request finished | %d | %s", id, lrw.Status(), dur)
}
