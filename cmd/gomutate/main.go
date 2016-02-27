package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/zabawaba99/gomutate"
)

var wd = getWD()

func getWD() string {
	wd, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatalf("Could not get working directory")
	}
	return wd
}

func main() {
	// parse flags
	opts := parseOptions()
	log.WithFields(log.Fields{"opts": opts}).Debug("Starting testing")

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}

	pkgs := opts.getPkgs()
	mutators := opts.getMutators()

	g := gomutate.New(wd)
	for _, pkg := range pkgs {
		g.Run(pkg, mutators)
	}

	g.GatherResults()
}
