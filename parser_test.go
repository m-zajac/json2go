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
			name:         "empty input",
			inputs:       []string{},
			expectedRepr: fmt.Sprintf("type %s interface{}", baseTypeName),
		},
		{
			name: "single simple input",
			inputs: []string{
				"1",
				"2",
			},
			expectedRepr: fmt.Sprintf("type %s int", baseTypeName),
		},
		{
			name: "single simple arrays",
			inputs: []string{
				"[1, 2.0]",
				"[3, 4]",
			},
			expectedRepr: fmt.Sprintf("type %s []int", baseTypeName),
		},
		{
			name: "single simple object",
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
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser()
			for _, v := range tc.inputs {
				if err := p.FeedBytes([]byte(v)); err != nil {
					t.Fatalf("feed error: %v", err)
				}
			}

			repr, err := p.String()
			if err != nil {
				t.Fatalf("parser.String error: %v", err)
			}
			expectedRepr := strings.TrimSpace(tc.expectedRepr)
			if repr != expectedRepr {
				t.Errorf("invalid repr.\nwant:\n%s\n\ngot:\n%s", expectedRepr, repr)
			}
		})
	}
}
