package main

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	"github.com/zabawaba99/gomutate/mutants"
)

var mutatorMapping = map[string]mutants.Mutator{
	"conditionals-boundary": &mutants.ConditionalsBoundary{},
	"negate-conditionals":   &mutants.NegateConditionals{},
}

type Options struct {
	Debug    bool     `short:"d" long:"debug" description:"Show debug information"`
	Mutator  []string `short:"m" long:"mutator" description:"The mutators to apply" default:"conditionals-boundary"`
	packages []string
}

func parseOptions() *Options {
	opts := &Options{}
	pkgs, err := flags.Parse(opts)
	if err != nil {
		log.Fatal("could not parse args")
	}

	opts.packages = pkgs
	return opts
}

func (opts *Options) getMutators() []mutants.Mutator {
	var mutators []mutants.Mutator
	for _, m := range opts.Mutator {
		mt, ok := mutatorMapping[m]
		if !ok {
			log.WithField("mutator", m).Warning("Unrecognized mutator")
			continue
		}
		mutators = append(mutators, mt)
	}
	return mutators
}

func (opts *Options) getPkgs() []string {
	pkgs := make([]string, len(opts.packages))
	var appendAll bool
	for i, pkg := range opts.packages {
		switch pkg {
		case "./...":
			appendAll = true
			fallthrough
		case ".":
			pkgs[i] = ""
		default:
			if !strings.HasPrefix(pkg, "./") {
				log.Fatalf("Package %q does not exist in current directory", pkg)
			}
			pkgs[i] = strings.TrimPrefix(pkg, "./")
		}
	}
	if appendAll {
		filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}

			if hasPrefix(info.Name(), "_", ".") {
				return nil
			}

			pkg := strings.TrimPrefix(path, wd)
			pkgs = append(pkgs, strings.TrimPrefix(pkg, "/"))
			return nil
		})
	}
	return dedup(pkgs)
}

func dedup(pkgs []string) []string {
	set := map[string]bool{}
	for _, pkg := range pkgs {
		set[pkg] = true
	}

	rtn := make([]string, len(set))
	counter := 0
	for v := range set {
		rtn[counter] = v
		counter++
	}
	return rtn
}

func hasPrefix(s string, prefixs ...string) bool {
	for _, prefix := range prefixs {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}
