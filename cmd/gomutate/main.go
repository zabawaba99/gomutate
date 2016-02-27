package main

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
	"github.com/zabawaba99/gomutate"
	"github.com/zabawaba99/gomutate/mutants"
)

type Options struct {
	Debug bool `short:"d" long:"debug" description:"Show debug information"`
}

func main() {
	var opts Options
	pkgs, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("could not parse args")
	}

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory %s", err)
	}

	var appendAll bool
	for i, pkg := range pkgs {
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
		pkgs = append(pkgs, findAllPackages(wd)...)
	}
	pkgs = dedup(pkgs)
	log.WithFields(log.Fields{"pkgs": pkgs}).Debug("User given packages")

	g := gomutate.New(wd)
	for _, pkg := range pkgs {
		g.Run(pkg, &mutants.ConditionalsBoundary{}, &mutants.NegateConditionals{})
	}

	g.GatherResults()
}

func findAllPackages(wd string) (pkgs []string) {
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		if hasPrefix(info.Name(), "_", ".") {
			return nil
		}

		pkgs = append(pkgs, strings.TrimPrefix(path, wd))
		return nil
	})
	return
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
