package json2go

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewSchemaItemFromNode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		node node
		want string
	}{
		{
			name: "bool",
			node: node{
				root:     true,
				t:        nodeTypeBool,
				nullable: false,
			},
			want: `
{
	"$schema": "http://json-schema.org/draft-07/schema",
	"$id": "http://example.com/example.json",
	"type": "boolean"
}
			`,
		},
		{
			name: "int",
			node: node{
				root:     true,
				t:        nodeTypeInt,
				nullable: false,
			},
			want: `
{
	"$schema": "http://json-schema.org/draft-07/schema",
	"$id": "http://example.com/example.json",
	"type": "integer"
}
			`,
		},
		{
			name: "float",
			node: node{
				root:     true,
				t:        nodeTypeFloat,
				nullable: false,
			},
			want: `
{
	"$schema": "http://json-schema.org/draft-07/schema",
	"$id": "http://example.com/example.json",
	"type": "number"
}
			`,
		},
		{
			name: "bool array",
			node: node{
				root:       true,
				t:          nodeTypeBool,
				nullable:   false,
				arrayLevel: 1,
			},
			want: `
{
	"$schema": "http://json-schema.org/draft-07/schema",
	"$id": "http://example.com/example.json",
	"type": "array"
}
			`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.want = strings.Trim(tc.want, " \t\n")
			item := newSchemaItemFromNode(tc.node)
			b, err := json.MarshalIndent(item, "", "\t")
			if err != nil {
				t.Fatalf("json marshall err: %v", err)
			}
			if v := string(b); v != tc.want {
				t.Errorf("newSchemaItemFromNode() =\n%s\n\twant\n%s", v, tc.want)
			}
		})
	}
}
