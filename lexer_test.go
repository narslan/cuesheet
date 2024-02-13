package cuesheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test lexer using file-driven tests.
func TestLexer(t *testing.T) {
	// Find the paths of all input files in the data directory.
	paths, err := filepath.Glob(filepath.Join("testdata", "*.lexer.input"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]

		// Each path turns into a test: the test name is the filename without the
		// extension.
		t.Run(testname, func(t *testing.T) {
			source, err := os.ReadFile(path)
			if err != nil {
				t.Fatal("error reading source file:", err)
			}

			// >>> This is the actual code under test.
			lexResult := lex_helper(t, string(source))
			var output strings.Builder
			lineSentinel := 1
			for _, item := range lexResult {
				if item.line != lineSentinel {
					output.WriteRune('\n')
					lineSentinel = lineSentinel + 1
				}
				output.WriteString(item.val)
				output.WriteRune(' ')
			}
			// <<<

			// Each input file is expected to have a "golden output" file, with the
			// same path except the .input extension is replaced by .golden
			goldenfile := filepath.Join("testdata", testname+".golden")
			want, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			if output.String() == string(want) {
				t.Errorf("\n==== got:\n%s\n==== want:\n%s\n", output.String(), want)
			}
		})
	}
}

func lex_helper(tb testing.TB, input string) []item {
	tb.Helper()
	result := make([]item, 0)
	l := lex(input)
	for {
		item := l.nextItem()
		result = append(result, item)
		if item.typ == itemEOF {
			break
		}
	}
	return result
}
