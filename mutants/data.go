package mutants

import (
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

const DataFileName = "gomutate.json"

type Data struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Original   string `json:"original"`
	Mutation   string `json:"mutation"`
	Type       string `json:"type"`
	Killed     bool   `json:"killed"`
}

func (md *Data) Load(path string) {
	f, err := os.Open(filepath.Join(path, DataFileName))
	if err != nil {
		log.WithError(err).Fatal("Could not open gomutate.json")
	}
	defer f.Close()

	json.NewDecoder(f).Decode(md)
}

func (md *Data) Save(path string) {
	f, err := os.Create(filepath.Join(path, DataFileName))
	if err != nil {
		log.WithError(err).Fatal("Could not create gomutate.json")
	}
	defer f.Close()

	json.NewEncoder(f).Encode(md)
}
