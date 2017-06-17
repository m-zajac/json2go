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
			name: "float arrays",
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
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			p := NewJSONParser()
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
