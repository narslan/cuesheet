package cuesheet

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/narslan/tree"
)

// TestParseSimpleStm tests simple one line commands. The FILE command has complex ones.
func TestParseSingleLines(t *testing.T) {
	f, err := os.Open("testdata/parser.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for lineno := 1; s.Scan(); lineno++ {
		line := s.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		t.Run("", func(t *testing.T) {

			l := strings.Split(line, "#")
			if len(l) != 2 {
				t.Error("testdata/parser.txt:", lineno, ": wrong field count")
				t.Fail()
			}

			nt := parse_helper(t, l[0])

			w := nt.Traverse()
			//exclude the root node
			w = w[1:]
			if len(w) == 0 {
				t.Fatal("parse failed...tree couldn't be built")
			}

			var b strings.Builder
			for i := range w {

				switch value := w[i].(type) {
				case node:
					b.WriteString(value.String())
					b.WriteRune(' ')
				}
			}
			e := strings.TrimSpace(l[1])       //expected
			g := strings.TrimSpace(b.String()) //got

			if e != g {

				t.Fatalf("expected %q got %q", e, g)
			}
		})

	}
}

// TestParseSimpleCueFiles parses common cue files  using file-driven tests.
func TestParseCueFiles(t *testing.T) {
	// Find the paths of all input files in the data directory.
	paths, err := filepath.Glob(filepath.Join("testdata", "*.parser.input"))
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
			parseResult := parse_helper(t, string(source))
			var output strings.Builder
			t.Log(parseResult)
			for _, item := range parseResult.Traverse()[1:] {
				switch v := item.(type) {
				case node:
					output.WriteString(v.Type.String())
					output.WriteRune(' ')
					output.WriteString(v.Value)
					output.WriteRune('\n')
				default:
					t.Fatal("unexpected type in parsed result")
				}
			}
			// <<<

			// Each input file is expected to have a "golden output" file, with the
			// same path except the .input extension is replaced by .golden
			goldenfile := filepath.Join("testdata", testname+".golden")
			want, err := os.ReadFile(goldenfile)
			if err != nil {
				t.Fatal("error reading golden file:", err)
			}

			if strings.TrimSpace(output.String()) != string(want) {
				t.Errorf("\n==== got:\n%q\n==== want:\n%q\n", output.String(), want)
			}
		})
	}
}

func parse_helper(tb testing.TB, input string) *tree.Tree {
	tb.Helper()

	parse := NewParser(input)
	tree, err := parse.Start()
	if err != nil {
		tb.Fatal(err)
	}

	return tree
}
