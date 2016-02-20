package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

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
	mutants, err := parseFile("sample/main.go")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Mutants %#v", mutants)

	a, err := newAST("sample/main.go")
	if err != nil {
		log.Fatal(err)
	}
	a.ApplyMutation(&BoundaryMutator{})
	// generate mutants

	// run tests

	// generate reports
}

type mutant struct {
	mutation    ast.Stmt
	mutationPos token.Pos
}

var mutationMap = map[token.Token]token.Token{
	//    <          >
	token.LSS: token.GTR,
	//    >          <
	token.GTR: token.LSS,
}

func parseFile(filename string) ([]mutant, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		return nil, err
	}

	var mutants []mutant
	// go through all declarations and find possible mutations
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			m := parseFunc(decl.(*ast.FuncDecl))
			mutants = append(mutants, m...)
		default:
			fmt.Printf("ParseFile: %T\n", decl)
		}
	}

	return mutants, nil
}

func parseFunc(fn *ast.FuncDecl) []mutant {
	var mutants []mutant
	for _, stmt := range fn.Body.List {
		switch stmt.(type) {
		case *ast.IfStmt:
			m := parseIfStmt(stmt.(*ast.IfStmt))
			mutants = append(mutants, m)
		default:
			fmt.Printf("ParseFunc: %T\n", stmt)
		}
	}
	return mutants
}

func parseIfStmt(ifStmt *ast.IfStmt) mutant {
	switch ifStmt.Cond.(type) {
	case *ast.BinaryExpr:
		bExpr := ifStmt.Cond.(*ast.BinaryExpr)
		mExpr := *bExpr
		mExpr.Op = mutationMap[bExpr.Op]
		mIfStmt := *ifStmt
		mIfStmt.Cond = &mExpr
		return mutant{
			mutation:    &mIfStmt,
			mutationPos: ifStmt.Pos(),
		}
	default:
		fmt.Printf("ParseIfStmty: %T\n", ifStmt)
	}

	return mutant{}
}
