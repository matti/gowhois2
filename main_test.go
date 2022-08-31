package gowhois

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVersion(t *testing.T) {
	if paths, err := filepath.Glob("tests/*"); err != nil {
		panic(err)
	} else {
		for _, path := range paths {
			if bytes, err := os.ReadFile(path); err != nil {
				panic(err)
			} else {
				Parse(string(bytes))
			}
		}
	}
}
