package gomutate

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/zabawaba99/gomutate/mutants"
)

type nodeVistor struct {
	count    int64
	mutation mutants.Mutator
	file     *token.File
	ast      *AST
}

func newNodeVisitor(a *AST, file *token.File, m mutants.Mutator) *nodeVistor {
	return &nodeVistor{ast: a, file: file, mutation: m}
}

func (v *nodeVistor) Visit(n ast.Node) ast.Visitor {
	v.ast.mtx.Lock()
	mutation, ok := v.mutation.Mutate(n)
	if !ok {
		v.ast.mtx.Unlock()
		return v
	}
	v.count++

	count := strconv.FormatInt(v.count, 10)
	filename := trimWD(v.file.Name())
	basedir := filepath.Join(mutationDir, v.mutation.Name(), filename+"."+count)
	if err := v.ast.write(basedir); err != nil {
		log.Fatalf("Could not create mutation file %s", err)
	}
	log.Debugf("Created mutation for %s", basedir)

	mutation.Reset()
	v.ast.mtx.Unlock()

	md := mutants.Data{
		Original:   mutation.OrgStmt,
		Mutation:   mutation.NewStmt,
		Type:       v.mutation.Name(),
		Filename:   v.file.Name(),
		LineNumber: v.file.Line(n.Pos()),
	}
	if err := md.Save(basedir); err != nil {
		fields := log.Fields{
			"error":    err,
			"mutation": basedir,
		}
		log.WithFields(fields).Warning("Could not save results for mutation")
	}

	return v
}
