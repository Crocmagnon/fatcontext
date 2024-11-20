package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
