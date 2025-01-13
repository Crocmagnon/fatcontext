package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}
	testdata := filepath.Join(wd, "testdata")

	t.Run("no func decl", func(t *testing.T) {
		an := analyzer.NewAnalyzer()
		analysistest.Run(t, testdata, an, "./common")
		analysistest.Run(t, testdata, an, "./no_structpointer")
	})

	t.Run("func decl", func(t *testing.T) {
		an := analyzer.NewAnalyzer()

		err := an.Flags.Set(analyzer.FlagCheckStructPointers, "true")
		if err != nil {
			t.Fatal(err)
		}

		analysistest.Run(t, testdata, an, "./common")
		analysistest.Run(t, testdata, an, "./structpointer")
	})
}
