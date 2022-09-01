package gowhois

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	if paths, err := filepath.Glob("inputs/*"); err != nil {
		panic(err)
	} else {
		for _, path := range paths {
			if bytes, err := os.ReadFile(path); err != nil {
				panic(err)
			} else {
				response := Parse(string(bytes))
				if bytes, err := os.ReadFile(
					filepath.Join("outputs", filepath.Base(path)),
				); err == nil {
					if string(bytes) != response.String()+"\n" {
						panic("differ")
					}
				} else {
					panic(err)
				}
			}
		}
	}
}
