package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name:     "osexit",
	Doc:      "prohibits using a direct os.Exit call",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	result := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	result.Preorder(nodeFilter, func(n ast.Node) {
		fnDecl := n.(*ast.FuncDecl)

		if fnDecl.Name.Name != "main" || pass.Pkg.Name() != "main" {
			return
		}

		ast.Inspect(fnDecl.Body, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if ident, ok := selExpr.X.(*ast.Ident); ok {
				if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
					pass.Reportf(callExpr.Pos(),
						"direct os.Exit call in main function is forbidden")
				}
			}
			return true
		})
	})

	return nil, nil
}
