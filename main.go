package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zabawaba99/gomutate/mutants"
)

const mutationDir = "_gomutate"

func init() {
	if err := os.RemoveAll(mutationDir); err != nil {
		fLog("Could not delete 'gomutate' directory %s\n", err)
	}

	if err := os.Mkdir(mutationDir, 0777); err != nil {
		fLog("Could not recreate 'gomutate' directory\n", err)
	}
}

func main() {
	// parse files
	a, err := newAST(wd)
	if err != nil {
		fLog("Could not read dir %s\n", err)
	}

	mutations := []mutants.Mutator{&mutants.NegateConditionals{}}
	for _, m := range mutations {
		// generate mutants
		a.ApplyMutation(m)

		// run tests
		runTests(m)
	}

	// generate reports
	aggregateResults()
}

func runTests(m mutants.Mutator) {
	mtpath := filepath.Join(mutationDir, m.Name())
	deviants, err := ioutil.ReadDir(mtpath)
	if err != nil {
		fLog("Could not find mutant directories %s", err)
	}

	for _, mt := range deviants {
		pkg := filepath.Join(mtpath, mt.Name())
		dLog("Running tests for %s", pkg)

		cmd := exec.Command("go", "test", "."+separator+pkg+separator+"...")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmd.Run()

		var md mutants.Data
		md.Load(pkg)
		md.Killed = !cmd.ProcessState.Success()
		md.Save(pkg)
	}
}

func aggregateResults() {
	results := []mutants.Data{}
	filepath.Walk(mutationDir, func(path string, info os.FileInfo, err error) error {
		if info.Name() != mutants.DataFileName {
			return nil
		}

		pkg := strings.TrimSuffix(path, info.Name())

		var result mutants.Data
		result.Load(pkg)
		results = append(results, result)

		return nil
	})

	f, err := os.Create(filepath.Join(mutationDir, "results.json"))
	if err != nil {
		fLog("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(results)
}
