package gomutate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
		log.WithField("mutation", m.Name()).Info("Generating mutation...")
		// generate mutations
		a.ApplyMutation(m)

		log.WithField("mutation", m.Name()).Info("Testing mutations...")
		// run tests
		g.runTests(pkg, m)
	}
}

func (g *Gomutate) runTests(pkg string, m mutants.Mutator) {
	mtpath := filepath.Join(mutationDir, m.Name(), pkg)
	deviants, err := ioutil.ReadDir(mtpath)
	if err != nil {
		log.WithField("pkg", mtpath).Warning("Could not find mutant directories")
		return
	}

	reg := regexp.MustCompile(`^.+\.go\.\d+$`)
	for _, mt := range deviants {
		if !mt.IsDir() {
			continue
		}

		if !reg.MatchString(mt.Name()) {
			// a subpackage
			continue
		}

		mutant := filepath.Join(mtpath, mt.Name())
		log.Debugf("Running tests for %s", mutant)

		cmd := exec.Command("go", "test", "."+separator+mutant+separator+"...")
		cmd.Run()

		var md mutants.Data
		md.Load(mutant)
		md.Killed = !cmd.ProcessState.Success()
		log.WithFields(log.Fields{"killed": md.Killed, "mutant": mutant}).Info("Result")
		md.Save(mutant)
	}
}

func (g *Gomutate) GatherResults() {
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
		log.WithError(err).Fatal("Could not create gomutate.json")
	}
	defer f.Close()

	json.NewEncoder(f).Encode(results)
}
