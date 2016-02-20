package main

import "os"

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

	a.ApplyMutation(&NegateConditionals{})
	// generate mutants

	// run tests

	// generate reports
}
