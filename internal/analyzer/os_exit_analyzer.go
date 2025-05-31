package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// OsExitAnalyzer prohibits the use of direct calls to os.Exit().
//
// This analyzer checks for direct invocations of os.Exit() within functions
// in the main package. Directly calling os.Exit() is discouraged because it
// bypasses normal error handling mechanisms and program lifecycle management.
// Instead, it's recommended to return an error from the main function,
// allowing the standard runtime handler to terminate the application properly.
var OsExitAnalyzer = &analysis.Analyzer{
	Name:     "osexit",
	Doc:      "prohibits using a direct os.Exit call",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer}, // Requires the inspect tool for analyzing AST
}

// run is executed when analysis runs. It traverses the abstract syntax tree
// looking for os.Exit() calls. If such a call is found, the analyzer reports a warning.
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
