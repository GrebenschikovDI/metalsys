// Package main is the entry point for staticlint command.
//
// To run the multichecker with custom analyzer:
//
//	go install ./cmd/staticlint
//	staticlint ./your-package
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := groupingAnalyzers(StaticAnalyzers(), standartAnalyzers, []*analysis.Analyzer{exitAnalyzer})
	multichecker.Main(
		analyzers...,
	)
}

func groupingAnalyzers(a ...[]*analysis.Analyzer) []*analysis.Analyzer {
	joined := make([]*analysis.Analyzer, 0, 100)
	for _, l := range a {
		joined = append(joined, l...)
	}
	return joined
}
