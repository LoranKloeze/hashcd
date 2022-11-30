package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lorankloeze/hashcd/log"

	"github.com/urfave/negroni"
)

type contextKey string

const ContextRequestIdKey contextKey = "requestId"

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

func Cors(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Expose-Headers", "X-Served-From")
	next(rw, r)
}

func RequestId(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := context.WithValue(r.Context(), ContextRequestIdKey, uuid.New())
	r = r.WithContext(ctx)
	next(rw, r)
}

func Log(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tStart := time.Now()

	id := r.Context().Value(ContextRequestIdKey)
	ctx := log.WithLogger(r.Context(), log.L.WithField("reqid", id))

	log.G(ctx).Infof("Request from %s [%s] %s", clientIp(r), r.Method, r.URL.Path)

	lrw := negroni.NewResponseWriter(rw)
	next(lrw, r)
	dur := time.Since(tStart)

	log.G(ctx).Infof("Request finished <%d - %s> in %s", lrw.Status(), http.StatusText(lrw.Status()), dur)
}
