package gomutate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

var (
	separator = string(filepath.Separator)
	wd        string
)

func init() {
	w, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory %s", err)
	}
	wd = w
}

func randomStr(size int) string {
	secs := time.Now().UnixNano()
	return strconv.FormatInt(secs, 10)
}

func trimWD(filename string) string {
	return strings.TrimPrefix(filename, wd+"/")
}
