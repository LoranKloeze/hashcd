package middleware

import (
	"bytes"
	"context"
	"regexp"

	log "github.com/sirupsen/logrus"

	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestCors(t *testing.T) {

	rw := httptest.NewRecorder()
	w := httptest.NewRequest("GET", "/", nil)

	Cors(rw, w, func(w http.ResponseWriter, r *http.Request) {
		testValues := [][]string{
			{"Access-Control-Allow-Origin", "*"},
			{"Access-Control-Expose-Headers", "X-Served-From"},
		}

		for _, tv := range testValues {
			if w.Header().Values(tv[0])[0] != tv[1] {
				t.Errorf("expected header %q to have value %q, got %q", tv[0], tv[1], w.Header().Values(tv[0])[0])
			}
		}
	})

}

func TestRequestId(t *testing.T) {

	rw := httptest.NewRecorder()
	w := httptest.NewRequest("GET", "/", nil)

	RequestId(rw, w, func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ContextRequestIdKey)
		if reflect.TypeOf(id).String() != "uuid.UUID" {
			t.Errorf("expected context key %q to have a value of type uuid.UUID, got %q", ContextRequestIdKey, reflect.TypeOf(id).String())
		}
	})

}

func TestLog(t *testing.T) {
	rw := httptest.NewRecorder()
	w := httptest.NewRequest("GET", "/", nil)
	w = w.WithContext(context.WithValue(w.Context(), ContextRequestIdKey, "myid"))
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Log(rw, w, func(w http.ResponseWriter, r *http.Request) {})

	for _, want := range []string{`myid`, `192.0.2.1`, `Request from`, `Request finished`} {
		re := regexp.MustCompile(want)
		if !re.MatchString(buf.String()) {
			t.Log(buf.String())
			t.Errorf("expected log to match %q", want)
		}
	}

}
