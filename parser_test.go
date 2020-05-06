package json2go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

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

	// Output: type Document struct {
	// 	Line	*struct {
	// 		End	XY	`json:"end"`
	// 		Start	XY	`json:"start"`
	// 	}	`json:"line,omitempty"`
	// 	Triangle	[]XY	`json:"triangle,omitempty"`
	// }
	// type XY struct {
	// 	X	float64	`json:"x"`
	// 	Y	float64	`json:"y"`
	// }
}

func ExampleJSONParser_FeedValue() {
	var v interface{}
	_ = json.Unmarshal([]byte(`{"line":{"start":{"x":12.1,"y":2.8},"end":{"x":12.1,"y":5.67}}}`), &v)

	parser := NewJSONParser("Document")
	parser.FeedValue(v)
	res := parser.String()
	fmt.Println(res)

	// Output: type Document struct {
	// 	Line struct {
	// 		End	struct {
	// 			X	float64	`json:"x"`
	// 			Y	float64	`json:"y"`
	// 		}	`json:"end"`
	// 		Start	struct {
	// 			X	float64	`json:"x"`
	// 			Y	float64	`json:"y"`
	// 		}	`json:"start"`
	// 	} `json:"line"`
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
			ExtractCommonTypes           bool `yaml:"extractCommonTypes"`
			StringPointersWhenKeyMissing bool `yaml:"stringPointersWhenKeyMissing"`
			SkipEmptyKeys                bool `yaml:"skipEmptyKeys"`
		} `yaml:"options"`
		Out string `yaml:"out"`
	}

	input, err := ioutil.ReadFile(inPath)
	require.NoError(t, err)

	output, err := ioutil.ReadFile(outPath)
	require.NoError(t, err)

	var tests []testDef
	err = yaml.Unmarshal(output, &tests)
	require.NoError(t, err)
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%s-%d", name, i), func(t *testing.T) {
			parserOpts := []JSONParserOpt{
				OptExtractCommonTypes(tc.Options.ExtractCommonTypes),
				OptStringPointersWhenKeyMissing(tc.Options.StringPointersWhenKeyMissing),
				OptSkipEmptyKeys(tc.Options.SkipEmptyKeys),
			}
			parser := NewJSONParser(baseTypeName, parserOpts...)
			err = parser.FeedBytes(input)
			require.NoError(t, err)

			got := normalizeStr(parser.String())
			want := normalizeStr(tc.Out)

			assert.Equal(t, want, got)
		})
	}
}

// normalizeStr trims string, replaces all tabs and space groups with single space, collapses multiple new lines into one.
func normalizeStr(v string) string {
	v = strings.TrimSpace(v)
	v = regexp.MustCompile(`[^\S\r\n]+`).ReplaceAllString(v, " ")
	v = regexp.MustCompile(`[\r\n]+`).ReplaceAllString(v, "\n")

	return v
}
