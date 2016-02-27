package gomutate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/zabawaba99/gomutate/mutants"
)

type AST struct {
	mtxs  map[string]sync.RWMutex
	fset  *token.FileSet
	pkgs  map[string]*ast.Package
	files map[string]*token.File
}

func newAST(filename string) (*AST, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	a := &AST{
		mtxs:  map[string]sync.RWMutex{},
		fset:  fset,
		pkgs:  pkgs,
		files: map[string]*token.File{},
	}
	fset.Iterate(func(f *token.File) bool {
		a.files[trimWD(f.Name())] = f
		return true
	})

	log.WithField("pkgs", fmt.Sprintf("%v", a.pkgs)).Debug("Parsed packages")

	return a, nil
}

type walkFunc func(name string, file *ast.File) error

func (a *AST) forEachFile(fn walkFunc) error {
	for _, pkg := range a.pkgs {
		for fname, file := range pkg.Files {
			fname = trimWD(fname)
			if err := fn(fname, file); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AST) ApplyMutation(m mutants.Mutator) {
	var wg sync.WaitGroup
	a.forEachFile(func(name string, file *ast.File) error {
		log.WithField("file", name).Debug("Visting file...")
		if strings.HasSuffix(name, "_test.go") {
			return nil
		}

		wg.Add(1)
		go func() {
			visitor := newNodeVisitor(a, a.files[name], m)
			ast.Walk(visitor, file)
			wg.Done()
		}()
		return nil
	})
	wg.Wait()
}

func (a *AST) write(basepath string) error {
	return a.forEachFile(func(name string, ast *ast.File) error {
		mtx := a.mtxs[name]
		mtx.RLock()
		defer mtx.RUnlock()

		filename := filepath.Join(basepath, name)
		if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		err = printer.Fprint(file, a.fset, ast)
		if err != nil {
			return err
		}

		return nil
	})
}
