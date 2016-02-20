package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const mutantDataName = "gomutate.json"

type MutantData struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Killed     bool   `json:"killed"`
}

func (md *MutantData) load(path string) {
	f, err := os.Open(filepath.Join(path, mutantDataName))
	if err != nil {
		fLog("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewDecoder(f).Decode(md)
}

func (md *MutantData) save(path string) {
	f, err := os.Create(filepath.Join(path, mutantDataName))
	if err != nil {
		fLog("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(md)
}
