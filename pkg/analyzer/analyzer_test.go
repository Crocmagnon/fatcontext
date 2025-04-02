package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc    string
		dir     string
		options map[string]string
	}{
		{
			desc: "no func decl",
			dir:  "common",
		},
		{
			desc: "no func decl",
			dir:  "no_structpointer",
		},
		{
			desc: "func decl",
			dir:  "common",
			options: map[string]string{
				analyzer.FlagCheckStructPointers: "true",
			},
		},
		{
			desc: "func decl",
			dir:  "structpointer",
			options: map[string]string{
				analyzer.FlagCheckStructPointers: "true",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc+"_"+test.dir, func(t *testing.T) {
			t.Parallel()

			anlzr := analyzer.NewAnalyzer()

			for k, v := range test.options {
				err := anlzr.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), anlzr, test.dir)
		})
	}
}

func TestAnalyzer_cgo(t *testing.T) {
	t.Parallel()

	a := analyzer.NewAnalyzer()

	analysistest.Run(t, analysistest.TestData(), a, "cgo")
}
