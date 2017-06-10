package jsontogo

import (
	"fmt"
	"testing"
)

func TestParserRepr(t *testing.T) {
	testCases := []struct {
		name         string
		inputs       []interface{}
		expectedRepr string
	}{
		{
			name:         "base type",
			inputs:       []interface{}{true, false, true},
			expectedRepr: fmt.Sprintf("type %s bool", baseTypeName),
		},
		{
			name: "slice type",
			inputs: []interface{}{
				[]interface{}{true, false},
				[]interface{}{false, true},
			},
			expectedRepr: fmt.Sprintf("type %s []bool", baseTypeName),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := NewParser()
			for _, v := range tc.inputs {
				p.Feed(v)
			}

			repr := p.String()
			if repr != tc.expectedRepr {
				t.Errorf("invalid repr.\nwant:\n%s\n\ngot:\n%s", tc.expectedRepr, repr)
			}
		})
	}
}
