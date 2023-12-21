package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var exitAnalyzer = &analysis.Analyzer{
	Name: "mainExitCheck",
	Doc:  "Check for direct os.Exit calls in main functions",
	Run:  runExitAnalyzer,
}

func runExitAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name == "main" {
			continue
		}
		inMain := false
		ast.Inspect(file, func(node ast.Node) bool {
			if isNameFunc(node, "main") {
				inMain = true
				return true
			}
			if inMain && isNameFunc(node, "Exit") {
				pass.Reportf(node.Pos(), "don't use exit in main")
			}

			return true
		})
	}
	return nil, nil
}

func isNameFunc(node ast.Node, name string) bool {
	if id, ok := node.(*ast.Ident); ok {
		if id.Name == name {
			return true
		}
	}
	return false
}
