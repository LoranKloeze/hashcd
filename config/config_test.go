package config

import (
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
