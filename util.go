package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fLog(msg string, args ...interface{}) {
	fmt.Printf(msg, args)
	os.Exit(1)
}

func randomStr(size int) string {
	secs := time.Now().UnixNano()
	return strconv.FormatInt(secs, 10)
}
