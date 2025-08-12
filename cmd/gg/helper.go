package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
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

var (
	green  = color.New(color.FgHiGreen).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	gray   = color.New(color.FgHiBlack).SprintFunc()
)

func logCreate(filename string) {
	fmt.Printf("%s %s\n", green("[CREATE]"), filename)
}

func logUpdate(filename string) {
	fmt.Printf("%s %s\n", yellow("[UPDATE]"), filename)
}

func logSkip(filename string) {
	fmt.Printf("%s %s\n", gray("[SKIP]"), filename)
}

func writeFileWithLog(filename string, content string) {
	if fileExists(filename) {
		oldData, err := os.ReadFile(filename)
		checkErr(err)
		if string(oldData) == content {
			logSkip(filename)
		} else {
			logUpdate(filename)
			checkErr(os.WriteFile(filename, []byte(content), 0o644))
		}
	} else {
		logCreate(filename)
		checkErr(ensureParentDir(filename))
		checkErr(os.WriteFile(filename, []byte(content), 0o644))
	}
}
