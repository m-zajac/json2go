package jsontogo

import (
	"fmt"
	"testing"
)

func TestFieldRepr(t *testing.T) {
	testCases := []struct {
		name           string
		startFieldName string
		startAsRoot    bool
		expands        []interface{}
		expectedRepr   string
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
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := newField(tc.startFieldName)
			f.root = tc.startAsRoot

			for _, v := range tc.expands {
				f.grow(v)
			}

			repr := f.repr()
			if repr != tc.expectedRepr {
				t.Errorf("invalid repr.\nwant:\n%s\n\ngot:\n%s", tc.expectedRepr, repr)
			}
		})
	}
}
