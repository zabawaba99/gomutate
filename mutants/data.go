package mutants

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const dataName = "gomutate.json"

type Data struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Killed     bool   `json:"killed"`
}

func (md *Data) Load(path string) {
	f, err := os.Open(filepath.Join(path, dataName))
	if err != nil {
		fLog("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewDecoder(f).Decode(md)
}

func (md *Data) Save(path string) {
	f, err := os.Create(filepath.Join(path, dataName))
	if err != nil {
		fLog("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(md)
}
