package gomutate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/zabawaba99/gomutate/mutants"
)

const mutationDir = "_gomutate"

func init() {
	if err := os.RemoveAll(mutationDir); err != nil {
		log.Fatalf("Could not delete '_gomutate' directory %s", err)
	}

	if err := os.Mkdir(mutationDir, 0777); err != nil {
		log.Fatalf("Could not recreate '_gomutate' directory %s", err)
	}
}

type Gomutate struct {
	wd string
}

func New(wd string) *Gomutate {
	return &Gomutate{wd: wd}
}

func (g *Gomutate) Run(pkg string, mutations ...mutants.Mutator) {
	// parse files
	a, err := newAST(filepath.Join(g.wd, pkg))
	if err != nil {
		log.Fatalf("Could not read dir %s", err)
	}

	for _, m := range mutations {
		log.Infof("Generating %s mutations", m.Name())
		// generate mutations
		a.ApplyMutation(m)

		log.Infof("Testing mutations")
		// run tests
		g.runTests(pkg, m)
	}

	// generate reports
	g.aggregateResults()
}

func (g *Gomutate) runTests(pkg string, m mutants.Mutator) {
	mtpath := filepath.Join(mutationDir, m.Name(), pkg)
	deviants, err := ioutil.ReadDir(mtpath)
	if err != nil {
		log.Debugf("Could not find mutant directories %s", err)
		return
	}

	for _, mt := range deviants {
		if !mt.IsDir() {
			continue
		}

		pkg := filepath.Join(mtpath, mt.Name())
		log.Debugf("Running tests for %s", pkg)

		cmd := exec.Command("go", "test", "."+separator+pkg+separator+"...")
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmd.Run()

		var md mutants.Data
		md.Load(pkg)
		md.Killed = !cmd.ProcessState.Success()
		log.Debugf("Killed %t", md.Killed)
		md.Save(pkg)
	}
}

func (g *Gomutate) aggregateResults() {
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
		log.Fatalf("Could not create gomutate.json %s", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(results)
}
