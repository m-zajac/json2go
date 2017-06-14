package json2go

import (
	"fmt"
	"strings"
	"testing"
)

func TestNodeRepr(t *testing.T) {
	testCases := []struct {
		name          string
		startNodeName string
		startAsRoot   bool
		expands       []interface{}
		expectedRepr  string
	}{
		// base types
		{
			name:         "bool",
			startAsRoot:  true,
			expands:      []interface{}{true},
			expectedRepr: fmt.Sprintf("type %s bool", baseTypeName),
		},
		{
			name:         "int",
			startAsRoot:  true,
			expands:      []interface{}{1},
			expectedRepr: fmt.Sprintf("type %s int", baseTypeName),
		},
		{
			name:         "float",
			startAsRoot:  true,
			expands:      []interface{}{1.1},
			expectedRepr: fmt.Sprintf("type %s float64", baseTypeName),
		},
		{
			name:         "string",
			startAsRoot:  true,
			expands:      []interface{}{"1.1"},
			expectedRepr: fmt.Sprintf("type %s string", baseTypeName),
		},
		{
			name:         "interface",
			startAsRoot:  true,
			expands:      []interface{}{true, 1},
			expectedRepr: fmt.Sprintf("type %s interface{}", baseTypeName),
		},

		// slice types
		{
			name:        "[]bool",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true},
			},
			expectedRepr: fmt.Sprintf("type %s []bool", baseTypeName),
		},
		{
			name:        "[]int",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1},
			},
			expectedRepr: fmt.Sprintf("type %s []int", baseTypeName),
		},
		{
			name:        "[]float",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1.1},
			},
			expectedRepr: fmt.Sprintf("type %s []float64", baseTypeName),
		},
		{
			name:        "[]string",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{"1.1"},
			},
			expectedRepr: fmt.Sprintf("type %s []string", baseTypeName),
		},
		{
			name:        "[]interface",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true, 1},
			},
			expectedRepr: fmt.Sprintf("type %s []interface{}", baseTypeName),
		},

		// objects
		{
			name:        "object with 1 attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": 1,
				},
				map[string]interface{}{
					"x": 3,
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X int `+"`json:\"x\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object with multiple attrs - sorted by name",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"c": 1,
					"b": 2.0,
					"a": "str",
				},
				map[string]interface{}{
					"c": 12.7,
					"b": true,
					"a": "str2",
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	A	string		`+"`json:\"a\"`"+`
	B	interface{}	`+"`json:\"b\"`"+`
	C	float64		`+"`json:\"c\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object with slice attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"slice": []interface{}{1, 2},
				},
				map[string]interface{}{
					"slice": []interface{}{3, 4, 5},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Slice []int `+"`json:\"slice\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object with nested object attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{12.5, 5},
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{6, 5.2},
					},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Level1 struct {
		Level2 []float64 `+"`json:\"level2\"`"+`
	} `+"`json:\"level1\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object 2 different attrs",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": 1,
				},
				map[string]interface{}{
					"y": "asd",
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	X	int	`+"`json:\"x\"`"+`
	Y	string	`+"`json:\"y\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "nested slice of objects attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{
							map[string]interface{}{
								"x": 1,
								"y": 2,
							},
							map[string]interface{}{
								"x": 3,
								"y": 4,
							},
						},
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{
							map[string]interface{}{
								"x": 345,
								"y": 45,
							},
							map[string]interface{}{
								"x": 6,
								"y": 45,
							},
						},
					},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Level1 struct {
		Level2 []struct {
			X	int	`+"`json:\"x\"`"+`
			Y	int	`+"`json:\"y\"`"+`
		} `+"`json:\"level2\"`"+`
	} `+"`json:\"level1\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "nested object with no data",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": map[string]interface{}{},
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": map[string]interface{}{},
					},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Level1 struct {
		Level2 interface{} `+"`json:\"level2\"`"+`
	} `+"`json:\"level1\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object attr - null then data",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": nil,
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": map[string]interface{}{
							"x": 1,
							"y": 2,
						},
					},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Level1 struct {
		Level2 struct {
			X	int	`+"`json:\"x\"`"+`
			Y	int	`+"`json:\"y\"`"+`
		} `+"`json:\"level2\"`"+`
	} `+"`json:\"level1\"`"+`
}
			`, baseTypeName),
		},
		{
			name:        "object attr - data then null",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": map[string]interface{}{
							"x": 1,
							"y": 2,
						},
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": nil,
					},
				},
			},
			expectedRepr: fmt.Sprintf(`
type %s struct {
	Level1 struct {
		Level2 struct {
			X	int	`+"`json:\"x\"`"+`
			Y	int	`+"`json:\"y\"`"+`
		} `+"`json:\"level2\"`"+`
	} `+"`json:\"level1\"`"+`
}
					`, baseTypeName),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newNode(tc.startNodeName)
			f.root = tc.startAsRoot

			for _, v := range tc.expands {
				f.grow(v)
			}

			repr, err := f.repr()
			if err != nil {
				t.Fatalf("got repr error.\nerror: %v\n repr:\n%s\n", err, repr)
			}

			expectedRepr := strings.TrimSpace(tc.expectedRepr)
			if repr != expectedRepr {
				t.Errorf("invalid repr.\nwant:\n%s\n\ngot:\n%s", expectedRepr, repr)
			}
		})
	}
}
