package gowhois

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	if paths, err := filepath.Glob("tests/*.txt"); err != nil {
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
