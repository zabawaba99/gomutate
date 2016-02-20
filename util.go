package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	separator = string(filepath.Separator)
	wd        string
)

func init() {
	w, err := os.Getwd()
	if err != nil {
		fLog("Could not get working directory %s\n", err)
	}
	wd = w
}

func dLog(msg string, args ...interface{}) {
	msg = fmt.Sprintf("DEBUG\t%s\n", msg)
	fmt.Printf(msg, args...)
}
func fLog(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}

func randomStr(size int) string {
	secs := time.Now().UnixNano()
	return strconv.FormatInt(secs, 10)
}
