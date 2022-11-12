package main

import (
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

func RandStringBytesMaskImprSrc(n int) string {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, (n+1)/2) // can be simplified to n/2 if n is always even

	if _, err := src.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)[:n]
}

func Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	maxMem := int64(10 << 20) // 10MB
	err := r.ParseMultipartForm(maxMem)
	if err != nil {
		log.Printf("Could not parse form: %s\n", err)
		http.Error(w, "Could not parse form data", http.StatusUnprocessableEntity)
		return
	}

	field, ok := r.MultipartForm.File["f"]
	if !ok {
		log.Println("Form field 'f' not found")
		http.Error(w, "Form field 'f' not found", http.StatusUnprocessableEntity)
		return
	}
	f, err := field[0].Open()
	if err != nil {
		log.Fatalf("Could not open file from form: %s\n", err)
	}
	defer f.Close()

	f1, err := os.Create("/home/loran/git/lab/mycdn/storage/myfile")
	if err != nil {
		log.Fatalf("Could not create file: %s\n", err)
	}
	defer f1.Close()

	io.Copy(f1, f)

}

func main() {
	router := httprouter.New()
	router.POST("/", Upload)
	// router.GET("/hello/:name", Hello)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
