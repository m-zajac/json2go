package json2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func ExampleNewJSONParser() {
	inputs := []string{
		`{"line":{"start":{"x":12.1,"y":2.8},"end":{"x":12.1,"y":5.67}}}`,
		`{"triangle":[{"x":2.34,"y":2.1}, {"x":45.1,"y":6.7}, {"x":4,"y":94.6}]}`,
	}

	parser := NewJSONParser(
		"Document",
		OptExtractCommonTypes(true),
	)

	for _, in := range inputs {
		_ = parser.FeedBytes([]byte(in))
	}

	res := parser.String()
	fmt.Println(res)

	// Output:
	// type Document struct {
	// 	Line *struct {
	// 		End   XY `json:"end"`
	// 		Start XY `json:"start"`
	// 	} `json:"line,omitempty"`
	// 	Triangle []XY `json:"triangle,omitempty"`
	// }
	// type XY struct {
	// 	X float64 `json:"x"`
	// 	Y float64 `json:"y"`
	// }

}

func ExampleJSONParser_FeedValue() {
	var v interface{}
	_ = json.Unmarshal([]byte(`{"line":{"start":{"x":12.1,"y":2.8},"end":{"x":12.1,"y":5.67}}}`), &v)

	parser := NewJSONParser("Document")
	parser.FeedValue(v)
	res := parser.String()
	fmt.Println(res)

	// Output:
	// type Document struct {
	// 	Line struct {
	// 		End   XY `json:"end"`
	// 		Start XY `json:"start"`
	// 	} `json:"line"`
	// }
	// type XY struct {
	// 	X float64 `json:"x"`
	// 	Y float64 `json:"y"`
	// }
}

// TestParser tests all cases from files in test/parser directory.
func TestParser(t *testing.T) {
	testfilesDir := "test/parser/"
	err := filepath.Walk(testfilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if ext != ".json" {
			return nil
		}
		outFile := strings.TrimSuffix(path, ext) + ".out.yaml"
		baseName := strings.TrimSuffix(path, ext)
		baseName = strings.TrimPrefix(baseName, testfilesDir)

		testFile(t, baseName, path, outFile)

		return nil
	})
	t.Log(err)
}

func testFile(t *testing.T, name, inPath, outPath string) {
	t.Helper()

	type testDef struct {
		Options struct {
			RootName                     string  `yaml:"rootName"`
			ExtractCommonTypes           bool    `yaml:"extractCommonTypes"`
			ExtractSimilarityThreshold   float64 `yaml:"extractSimilarityThreshold"`
			ExtractMinSubsetSize         int     `yaml:"extractMinSubsetSize"`
			ExtractMinSubsetOccurrences  int     `yaml:"extractMinSubsetOccurrences"`
			ExtractMinAddedFields        int     `yaml:"extractMinAddedFields"`
			StringPointersWhenKeyMissing bool    `yaml:"stringPointersWhenKeyMissing"`
			SkipEmptyKeys                bool    `yaml:"skipEmptyKeys"`
			MakeMaps                     bool    `yaml:"makeMaps"`
			MakeMapsWhenMinAttributes    uint    `yaml:"makeMapsWhenMinAttributes"`
			TimeAsStr                    bool    `yaml:"timeAsStr"`
			SkipGoRun                    bool    `yaml:"skipGoRun"`
		} `yaml:"options"`
		Out string `yaml:"out"`
	}

	input, err := os.ReadFile(inPath)
	require.NoError(t, err)

	output, err := os.ReadFile(outPath)
	require.NoError(t, err)

	var tests []testDef
	err = yaml.Unmarshal(output, &tests)
	require.NoError(t, err)
	for i, tc := range tests {
		tn := fmt.Sprintf("%s-%d", name, i)
		t.Run(tn, func(t *testing.T) {
			parserOpts := []JSONParserOpt{
				OptExtractCommonTypes(tc.Options.ExtractCommonTypes),
				OptStringPointersWhenKeyMissing(tc.Options.StringPointersWhenKeyMissing),
				OptSkipEmptyKeys(tc.Options.SkipEmptyKeys),
				OptMakeMaps(tc.Options.MakeMaps, tc.Options.MakeMapsWhenMinAttributes),
				OptTimeAsString(tc.Options.TimeAsStr),
			}
			if tc.Options.ExtractSimilarityThreshold > 0 {
				parserOpts = append(parserOpts, OptExtractHeuristics(
					tc.Options.ExtractSimilarityThreshold,
					tc.Options.ExtractMinSubsetSize,
					tc.Options.ExtractMinSubsetOccurrences,
					tc.Options.ExtractMinAddedFields,
				))
			}
			rootName := tc.Options.RootName
			if rootName == "" {
				rootName = baseTypeName
			}
			parser := NewJSONParser(rootName, parserOpts...)
			err = parser.FeedBytes(input)
			require.NoError(t, err)

			// Test .String() output 2 times to check if .String() doesn't change internal parser state.
			for i := 0; i < 2; i++ {
				parserOutput := parser.String()
				got := strings.TrimSpace(parserOutput)
				want := strings.TrimSpace(tc.Out)
				assert.Equal(t, want, got)
			}

			if !tc.Options.SkipGoRun {
				testGeneratedType(t, parser, input)
			}
		})
	}
}

// testGeneratedType unmarshals test data to generated type, then marshals it again and compares generated output to original data.
func testGeneratedType(t *testing.T, parser *JSONParser, data []byte) {
	err := parser.FeedBytes(data)
	require.NoError(t, err)

	parserOutput := parser.String()

	filename := makeTypeTestGoFile(t, parserOutput, parser.rootNode.name)

	runCmd := exec.Command("go", "run", filename)
	runCmd.Stdin = bytes.NewBuffer(data)
	out, err := runCmd.CombinedOutput()
	require.NoError(t, err, "running go code: %v, %s", err, out)

	// unmarshal input data and test output data to generic type, then compare
	var valIn, valOut interface{}
	err = json.Unmarshal(data, &valIn)
	require.NoError(t, err, "unmarshaling input data: %v", err)
	err = json.Unmarshal([]byte(out), &valOut)
	require.NoError(t, err, "unmarshaling output data: %v", err)
	if !compareIgnoringNilKeys(t, valIn, valOut) {
		t.Logf("got different value after marshal/unmarshal:\n%#+v\n%#+v\n\n%s", valIn, valOut, parserOutput)
	}
}

func makeTypeTestGoFile(t *testing.T, parserOutput, rootName string) string {
	testTemplate := `
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

{{.Type}}

func main() {
	var _ time.Time
	var doc {{.RootName}}

	jd := json.NewDecoder(os.Stdin)
	if err := jd.Decode(&doc); err != nil {
		fmt.Printf("json decoding error: %v\n", err)
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(doc)
}
`

	filename := path.Join(t.TempDir(), "main.go")
	f, err := os.Create(filename)
	require.NoError(t, err, "creating tmp file: %v", err)
	defer f.Close()

	tmpl, err := template.New("test").Parse(testTemplate)
	require.NoError(t, err, "parsing test code template: %v", err)
	err = tmpl.Execute(f, map[string]interface{}{
		"Type":     parserOutput,
		"RootName": rootName,
	})
	require.NoError(t, err, "executing test template: %v", err)

	return filename
}

// normalizeStr trims string, replaces all tabs and space groups with single space, collapses multiple new lines into one.
func normalizeStr(v string) string {
	v = strings.TrimSpace(v)
	v = regexp.MustCompile(`[^\S\r\n]+`).ReplaceAllString(v, " ")
	v = regexp.MustCompile(`[\r\n]+`).ReplaceAllString(v, "\n")

	return v
}

func compareIgnoringNilKeys(t *testing.T, a, b interface{}) bool {
	if ma, ok := a.(map[string]interface{}); ok {
		if mb, ok := b.(map[string]interface{}); ok {
			return compareMaps(t, ma, mb)
		}
	}
	if sa, ok := a.([]interface{}); ok {
		if sb, ok := b.([]interface{}); ok {
			return compareSlices(t, sa, sb)
		}
	}

	return assert.Equal(t, a, b)
}

func compareMaps(t *testing.T, a, b map[string]interface{}) bool {
	checkedKeys := make(map[string]bool)
	for k, v := range a {
		checkedKeys[k] = true
		if v == nil && b[k] != nil {
			t.Errorf("key '%s' is missing in one of the maps", k)
		} else if !compareIgnoringNilKeys(t, v, b[k]) {
			return false
		}
	}
	for k, v := range b {
		if checkedKeys[k] {
			continue
		}
		if v != nil {
			t.Errorf("key '%s' is missing in one of the maps", k)
		}
	}

	return true
}

func compareSlices(t *testing.T, a, b []interface{}) bool {
	if !assert.Equal(t, len(a), len(b)) {
		return false
	}
	for i := range a {
		if !compareIgnoringNilKeys(t, a[i], b[i]) {
			return false
		}
	}
	return true
}
