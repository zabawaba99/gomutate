package main

import "os"

var mutationDir = "_gomutate"

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
	wd, err := os.Getwd()
	if err != nil {
		fLog("Could not get working directory %s\n", err)
	}

	a, err := newAST(wd)
	if err != nil {
		fLog("Could not read dir %s\n", err)
	}

	a.ApplyMutation(&NegateConditionals{})
	// generate mutants

	// run tests

	// generate reports
}
