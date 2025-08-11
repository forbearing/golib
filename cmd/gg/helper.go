package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ensureParentDir(filename string) error {
	dir := filepath.Dir(filename)

	if _, err := os.Stat(dir); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	} else {
		return err
	}
}
