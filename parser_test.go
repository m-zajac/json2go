package json2go

import (
	"fmt"
	"strings"
	"testing"
)

func TestParserRepr(t *testing.T) {
	testCases := []struct {
		name         string
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
			name: "int arrays (should produce just int type)",
			inputs: []string{
				"[1, 2.0]",
				"[3, 4]",
			},
			expectedRepr: fmt.Sprintf("type %s int", baseTypeName),
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
			name: "simple object, one attr nullable",
			inputs: []string{
				`{"x": true, "y": "str"}`,
				`{"x": false}`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	*string	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
		{
			name: "array of objects, one attr nullable",
			inputs: []string{
				`[{"x": true, "y": "str"}, {"x": false}]`,
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	bool	`+"`json:\"x\"`"+`
	Y	*string	`+"`json:\"y,omitempty\"`"+`
}
			`, baseTypeName),
		},
		{
			name: "object, with nested object nullable slice",
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
			name: "object, attr with object slice with nullable attributes",
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
		A	*string	`+"`json:\"a,omitempty\"`"+`
		B	*string	`+"`json:\"b,omitempty\"`"+`
	}	`+"`json:\"y,omitempty\"`"+`
}
					`, baseTypeName),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			p := NewJSONParser(baseTypeName)
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
