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
		parser.FeedBytes([]byte(in))
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
	json.Unmarshal([]byte(`{"line":{"start":{"x":12.1,"y":2.8},"end":{"x":12.1,"y":5.67}}}`), &v)

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

func TestParserRepr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		opts         options
		inputs       []string
		expectedRepr string
	}{
		{
			name:         "empty",
			inputs:       []string{},
			expectedRepr: fmt.Sprintf("type %s interface{}", baseTypeName),
		},
		{
			name: "int",
			inputs: []string{
				"1",
				"2",
			},
			expectedRepr: fmt.Sprintf("type %s int", baseTypeName),
		},
		{
			name: "int arrays",
			inputs: []string{
				"[1, 2.0]",
				"[3, 4]",
			},
			expectedRepr: fmt.Sprintf("type %s []int", baseTypeName),
		},
		{
			name: "simple object",
			inputs: []string{
				`{"x": true, "y": "str"}`,
				`{"x": false, "y": "str2"}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	string	`+"`json:\"y\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "simple object, one attr not required",
			inputs: []string{
				`{"x": true, "y": "str"}`,
				`{"x": false}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	string	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "simple object with string attr nullable",
			inputs: []string{
				`{"x": "ok"}`,
				`{"x": null}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X *string `+"`json:\"x\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "simple object with string attr not required",
			inputs: []string{
				`{"x": "ok"}`,
				`{}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X string `+"`json:\"x,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "simple object with attr with only null values",
			inputs: []string{
				`{"x": null}`,
				`{"x": null}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X interface{} `+"`json:\"x\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "array of objects, one attr not required",
			inputs: []string{
				`[{"x": true, "y": "str"}, {"x": false}]`,
			},
			expectedRepr: fmt.Sprintf(`
type %s []struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	string	`+"`json:\"y,omitempty\"`"+`
}
			`, baseTypeName),
		},
		{
			name: "object, with nested object not required slice",
			inputs: []string{
				`{
					"x": true,
					"y": {
						"p": 1
					}
				}`,
				`{
					"x": true,
					"y": {
						"p": 1,
						"z": [
							{"a": 1, "b": 2},
							{"a": 1, "b": 2, "c": 1.23}
						]
					}
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	struct {
		P	int	`+"`json:\"p\"`"+`
		Z	[]struct {
			A	int		`+"`json:\"a\"`"+`
			B	int		`+"`json:\"b\"`"+`
			C	*float64	`+"`json:\"c,omitempty\"`"+`
		}	`+"`json:\"z,omitempty\"`"+`
	}	`+"`json:\"y\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "object, attr with object slice with not required attributes",
			inputs: []string{
				`{
					"x": true
				}`,
				`{
					"x": false,
					"y": [
						{"a": "yes"},
						{"b": "no"}
					]
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	[]struct {
		A	string	`+"`json:\"a,omitempty\"`"+`
		B	string	`+"`json:\"b,omitempty\"`"+`
	}	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "nested array of ints",
			inputs: []string{
				`[[1, 2, 3]]`,
				`[[4, 5, 6]]`,
			},
			expectedRepr: fmt.Sprintf("type %s [][]int", baseTypeName),
		},
		{
			name: "nested array of objects",
			inputs: []string{
				`[[{"x": false}, {"x": true}]]`,
				`[[{"x": true, "y": true}, {"x": false, "z": true}]]`,
			},
			expectedRepr: fmt.Sprintf(`
type %s [][]struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	*bool	`+"`json:\"y,omitempty\"`"+`
	Z	*bool	`+"`json:\"z,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "don't extract common type",
			inputs: []string{
				`{
					"x": {"z": 1}
				}`,
				`{
					"y": {"z": 2}
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	*struct {
		Z int `+"`json:\"z\"`"+`
	}	`+"`json:\"x,omitempty\"`"+`
	Y	*struct {
		Z int `+"`json:\"z\"`"+`
	}	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "extract common type",
			opts: options{
				extractCommonTypes: true,
			},
			inputs: []string{
				`{
					"x": {"z": 1}
				}`,
				`{
					"y": {"z": 2}
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	*Z	`+"`json:\"x,omitempty\"`"+`
	Y	*Z	`+"`json:\"y,omitempty\"`"+`
}
type Z struct {
	Z int `+"`json:\"z\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "no pointer strings",
			inputs: []string{
				`{
					"x": "ok"
				}`,
				`{
					"y": "ok"
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	string	`+"`json:\"x,omitempty\"`"+`
	Y	string	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "pointer strings",
			opts: options{
				stringPointersWhenKeyMissing: true,
			},
			inputs: []string{
				`{
					"x": "ok"
				}`,
				`{
					"y": "ok"
				}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	*string	`+"`json:\"x,omitempty\"`"+`
	Y	*string	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			p := NewJSONParser(baseTypeName)
			p.opts = tc.opts
			for _, v := range tc.inputs {
				if err := p.FeedBytes([]byte(v)); err != nil {
					t.Fatalf("feed error: %v", err)
				}
			}

			repr := p.String()
			expectedRepr := strings.TrimSpace(tc.expectedRepr)
			if repr != expectedRepr {
				t.Errorf("invalid repr.\nwant:\n%s\n\ngot:\n%s", expectedRepr, repr)
			}
		})
	}
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
