package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

	mutations := []Mutator{&NegateConditionals{}}
	for _, m := range mutations {
		// generate mutants
		a.ApplyMutation(m)

		// run tests
		runTests(m)
	}

	// generate reports
}

func runTests(m Mutator) {
	mtpath := filepath.Join(mutationDir, m.Name())
	mutants, err := ioutil.ReadDir(mtpath)
	if err != nil {
		fLog("Could not find mutant directories %s", err)
	}

	for _, mt := range mutants {
		pkg := filepath.Join(mtpath, mt.Name())
		dLog("Running tests for %s", pkg)

		cmd := exec.Command("go", "test", "."+separator+pkg+separator+"...")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmd.Run()

		var md MutantData
		md.load(pkg)
		md.Killed = !cmd.ProcessState.Success()
		md.save(pkg)
	}
}
