package server

import "os"

type Config struct {
	storageDir string
}

var config = Config{
	storageDir: os.Getenv("HASHCD_STORAGE"),
}
