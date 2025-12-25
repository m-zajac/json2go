package json2go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeSimilarity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		n1       *node
		n2       *node
		expected float64
	}{
		{
			name: "identical simple nodes",
			n1: &node{
				t: nodeTypeInt,
			},
			n2: &node{
				t: nodeTypeInt,
			},
			expected: 1.0,
		},
		{
			name: "different base types",
			n1: &node{
				t: nodeTypeInt,
			},
			n2: &node{
				t: nodeTypeString,
			},
			expected: 0.0,
		},
		{
			name: "identical objects",
			n1: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
					{key: "b", t: nodeTypeString},
				},
			},
			n2: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
					{key: "b", t: nodeTypeString},
				},
			},
			expected: 1.0,
		},
		{
			name: "partially similar objects",
			n1: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
					{key: "b", t: nodeTypeString},
				},
			},
			n2: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
					{key: "c", t: nodeTypeBool},
				},
			},
			// intersection: {"a|int"} (1)
			// union: {"a|int", "b|string", "c|bool"} (3)
			// expected: 1/3 = 0.333...
			expected: 1.0 / 3.0,
		},
		{
			name: "objects with same keys but different types",
			n1: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
				},
			},
			n2: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeString},
				},
			},
			// intersection: 0
			// union: {"a|int", "a|string"} (2)
			expected: 0.0,
		},
		{
			name: "one object empty",
			n1: &node{
				t: nodeTypeObject,
				children: []*node{
					{key: "a", t: nodeTypeInt},
				},
			},
			n2: &node{
				t: nodeTypeObject,
			},
			expected: 0.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.n1.similarity(tc.n2)
			assert.InDelta(t, tc.expected, actual, 0.0001)
		})
	}
}
