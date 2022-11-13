package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

func GenMD5Hash(f io.Reader) string {
	hash := md5.New()
	_, err := io.Copy(hash, f)
	if err != nil {
		log.Fatalf("Could not calculate MD5 hash %s\n", err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func InitDirectories(hash string) (string, error) {
	t := hash[0 : len(hash)-2]

	re := regexp.MustCompile(`..`)
	p := "/home/loran/git/lab/mycdn/storage/"
	r := re.FindAllString(t, -1)
	p += strings.Join(r, "/")
	err := os.MkdirAll(p, 0755)
	if err != nil {
		log.Errorf("Could not create directory storage tree: %s", err)
		return "", err
	}
	log.Debugf("Initialized directory storage tree '%s'", p)
	return p, nil
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
		log.Error("Could not open file from form: %s\n", err)
		http.Error(w, "Could not open file from form", http.StatusUnprocessableEntity)
		return
	}
	defer f.Close()

	hash := GenMD5Hash(f)
	path, err := InitDirectories(hash)
	if err != nil {
		log.Fatalf("Could not create directory storage tree %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f.Seek(0, 0)
	f1, err := os.Create(fmt.Sprintf("%s/%s", path, hash))
	if err != nil {
		log.Fatalf("Could not create file: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f1.Close()

	io.Copy(f1, f)
	w.WriteHeader(http.StatusCreated)
}

func main() {
	log.SetLevel(log.DebugLevel)
	router := httprouter.New()
	router.POST("/", Upload)
	// router.GET("/hello/:name", Hello)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
