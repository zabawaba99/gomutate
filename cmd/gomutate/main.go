package main

import (
	"log"
	"os"

	"github.com/zabawaba99/gomutate"
	"github.com/zabawaba99/gomutate/mutants"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory %s\n", err)
	}
	g := gomutate.New(wd)
	g.Run(&mutants.ConditionalsBoundary{}, &mutants.NegateConditionals{})
}
