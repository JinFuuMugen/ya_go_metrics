package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer reports panic usages and log.Fatal*/os.Exit outside main.
var Analyzer = &analysis.Analyzer{
	Name:     "forbidencalls",
	Doc:      "eports panic usages and log.Fatal*/os.Exit outside main",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if !push {
			return true
		}

		call := n.(*ast.CallExpr)

		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
			pass.Reportf(call.Lparen, "avoid using panic")
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		insideMain := isInsideMainMain(pass, stack)

		if pkgIdent.Name == "log" && (sel.Sel.Name == "Fatal" || sel.Sel.Name == "Fatalf" || sel.Sel.Name == "Fatalln") {
			if !insideMain {
				pass.Reportf(call.Lparen, "log.%s is not allowed outside main.main", sel.Sel.Name)
			}
			return true
		}

		if pkgIdent.Name == "os" && sel.Sel.Name == "Exit" {
			if !insideMain {
				pass.Reportf(call.Lparen, "os.Exit is not allowed outside main.main")
			}
			return true
		}

		return true
	})

	return nil, nil
}

func isInsideMainMain(pass *analysis.Pass, stack []ast.Node) bool {
	if pass.Pkg == nil || pass.Pkg.Name() != "main" {
		return false
	}
	for i := len(stack) - 1; i >= 0; i-- {
		if fd, ok := stack[i].(*ast.FuncDecl); ok && fd.Name != nil && fd.Name.Name == "main" {
			return true
		}
	}
	return false
}
