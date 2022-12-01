package config

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDefaultConfig(t *testing.T) {
	Load()

	type Test struct {
		name   string
		values []interface{}
	}

	tests := []Test{
		{name: "CacheSize", values: []interface{}{C.CacheSize, int64(128)}},
		{name: "CacheItemSize", values: []interface{}{C.CacheItemSize, int64(2)}},
		{name: "LogLevel", values: []interface{}{C.LogLevel, logrus.InfoLevel}},
		{name: "StorageDir", values: []interface{}{C.StorageDir, ""}},
		{name: "ListenAddr", values: []interface{}{C.ListenAddr, "127.0.0.1:8080"}},
	}

	for _, test := range tests {
		if test.values[0] != test.values[1] {
			t.Errorf("Default configuration for %q: expected %v, got %v", test.name, test.values[0], test.values[1])
		}
	}

}

func TestMissingStorageDir(t *testing.T) {

	err := Load()
	if err == nil { // yes, this is err == nil
		t.Errorf("Expected error, got nil")
	}

}

func TestLoadWithEnvironment(t *testing.T) {

	type Test struct {
		name   string
		values []interface{}
	}

	// set environment variables
	wd, _ := os.Getwd() // Safe to use since nothjing is changed to it

	os.Setenv("HASHCD_CACHE_SIZE", "256")
	os.Setenv("HASHCD_CACHE_ITEM_SIZE", "4")
	os.Setenv("HASHCD_LOGLEVEL", "debug")
	os.Setenv("HASHCD_STORAGE_DIR", wd)
	os.Setenv("HASHCD_LISTEN_ADDR", "127.0.0.01:8123")

	err := Load()
	if err != nil {
		t.Errorf("Expected no error loading configuration, got: %v", err)
	}

	tests := []Test{
		{name: "CacheSize", values: []interface{}{C.CacheSize, int64(256)}},
		{name: "CacheItemSize", values: []interface{}{C.CacheItemSize, int64(4)}},
		{name: "LogLevel", values: []interface{}{C.LogLevel, logrus.DebugLevel}},
		{name: "StorageDir", values: []interface{}{C.StorageDir, wd}},
		{name: "ListenAddr", values: []interface{}{C.ListenAddr, "127.0.0.01:8123"}},
	}

	for _, test := range tests {
		if test.values[0] != test.values[1] {
			t.Errorf("Configuration from environment for %q: expected %v, got %v", test.name, test.values[1], test.values[0])
		}
	}

}
