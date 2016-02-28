package mutants

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func (md *Data) Load(path string) error {
	f, err := os.Open(filepath.Join(path, DataFileName))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(md)
}

func (md *Data) Save(path string) error {
	f, err := os.Create(filepath.Join(path, DataFileName))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(md)
}
