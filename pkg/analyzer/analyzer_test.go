package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
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
			desc: "no func decl",
			dir:  "cgo",
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

			a := analyzer.NewAnalyzer()

			for k, v := range test.options {
				err := a.Flags.Set(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.Run(t, analysistest.TestData(), a, test.dir)
		})
	}
}
