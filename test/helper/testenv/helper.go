package testenv

import (
	"os"
	"path/filepath"
)

func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get current working directory: " + err.Error())
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod not found)")
		}

		dir = parent
	}
}
