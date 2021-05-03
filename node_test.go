package json2go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONNodeCompare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		n1            *node
		n2            *node
		expectedEqual bool
	}{
		{
			name: "name not equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
			},
			n2: &node{
				key: "n2",
				t:   nodeTypeBool,
			},
			expectedEqual: false,
		},
		{
			name: "type not equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
			},
			n2: &node{
				key: "n",
				t:   nodeTypeInt,
			},
			expectedEqual: false,
		},
		{
			name: "required not equal",
			n1: &node{
				key:      "n",
				t:        nodeTypeBool,
				required: true,
			},
			n2: &node{
				key:      "n",
				t:        nodeTypeBool,
				required: false,
			},
			expectedEqual: false,
		},
		{
			name: "nullable not equal",
			n1: &node{
				key:      "n",
				t:        nodeTypeBool,
				nullable: true,
			},
			n2: &node{
				key:      "n",
				t:        nodeTypeBool,
				nullable: false,
			},
			expectedEqual: false,
		},
		{
			name: "simple equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
			},
			n2: &node{
				key: "n",
				t:   nodeTypeBool,
			},
			expectedEqual: true,
		},
		{
			name: "complex equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
				},
			},
			n2: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
				},
			},
			expectedEqual: true,
		},
		{
			name: "complex - child num not equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
				},
			},
			n2: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
				},
			},
			expectedEqual: false,
		},
		{
			name: "complex - child type not equal",
			n1: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeFloat,
						arrayLevel: 1,
					},
				},
			},
			n2: &node{
				key: "n",
				t:   nodeTypeBool,
				children: []*node{
					{
						key:        "n1",
						t:          nodeTypeInterface,
						arrayLevel: 1,
					},
				},
			},
			expectedEqual: false,
		},
		{
			name: "array level equal",
			n1: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 5,
			},
			n2: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 5,
			},
			expectedEqual: true,
		},
		{
			name: "array level not equal #1",
			n1: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 0,
			},
			n2: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 1,
			},
			expectedEqual: false,
		},
		{
			name: "array level not equal #2",
			n1: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 1,
			},
			n2: &node{
				key:        "n",
				t:          nodeTypeBool,
				arrayLevel: 2,
			},
			expectedEqual: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			eq := tc.n1.compare(tc.n2)
			if eq != tc.expectedEqual {
				t.Errorf("want %t, got %t", tc.expectedEqual, eq)
			}
		})
	}
}

func TestJSONNodeRepr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		startAsRoot bool
		expands     []interface{}
		expected    *node
	}{
		// base types
		{
			name:        "bool",
			startAsRoot: true,
			expands:     []interface{}{true},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeBool,
				nullable: false,
				required: true,
			},
		},
		{
			name:        "int",
			startAsRoot: true,
			expands:     []interface{}{1},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeInt,
				nullable: false,
				required: true,
			},
		},
		{
			name:        "float",
			startAsRoot: true,
			expands:     []interface{}{1.1},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeFloat,
				nullable: false,
				required: true,
			},
		},
		{
			name:        "string",
			startAsRoot: true,
			expands:     []interface{}{"1.1"},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeString,
				nullable: false,
				required: true,
			},
		},
		{
			name:        "interface",
			startAsRoot: true,
			expands:     []interface{}{true, 1},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeInterface,
				nullable: false,
				required: true,
			},
		},

		// slice types
		{
			name:        "[]bool",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeBool,
				arrayLevel: 1,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[]int",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeInt,
				arrayLevel: 1,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[]float",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1.1},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeFloat,
				arrayLevel: 1,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[]string",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{"1.1"},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeString,
				arrayLevel: 1,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[]interface",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true, 1},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeInterface,
				arrayLevel: 1,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[][]bool",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{
					[]interface{}{true},
				},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeBool,
				arrayLevel: 2,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[][][]bool",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{
					[]interface{}{
						[]interface{}{true},
					},
				},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeBool,
				arrayLevel: 3,
				nullable:   false,
				required:   true,
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeInt,
						nullable: false,
						required: true,
					},
				},
			},
		},
		{
			name:        "object with multiple attrs - sorted by name",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"c": 1,
					"a": "str",
					"b": 2.0,
				},
				map[string]interface{}{
					"a": "str2",
					"c": 12.7,
					"b": true,
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "a",
						t:        nodeTypeString,
						nullable: false,
						required: true,
					},
					{
						key:      "b",
						t:        nodeTypeInterface,
						nullable: false,
						required: true,
					},
					{
						key:      "c",
						t:        nodeTypeFloat,
						nullable: false,
						required: true,
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:        "slice",
						t:          nodeTypeInt,
						arrayLevel: 1,
						nullable:   false,
						required:   true,
					},
				},
			},
		},
		{
			name:        "object with required string attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": "stringx1",
				},
				map[string]interface{}{
					"x": "stringx2",
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeString,
						nullable: false,
						required: true,
					},
				},
			},
		},
		{
			name:        "object with not required string attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": "stringx1",
				},
				map[string]interface{}{
					"y": "stringy1",
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeString,
						nullable: false,
						required: false,
					},
					{
						key:      "y",
						t:        nodeTypeString,
						nullable: false,
						required: false,
					},
				},
			},
		},
		{
			name:        "object with nullable attr",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": "stringx1",
				},
				map[string]interface{}{
					"x": nil,
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeString,
						nullable: true,
						required: true,
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeFloat,
								arrayLevel: 1,
								nullable:   false,
								required:   true,
							},
						},
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeInt,
						nullable: false,
						required: false,
					},
					{
						key:      "y",
						t:        nodeTypeString,
						nullable: false,
						required: false,
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeObject,
								arrayLevel: 1,
								nullable:   false,
								required:   true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
								},
							},
						},
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeObject,
								nullable: false,
								required: true,
							},
						},
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeObject,
								nullable: true,
								required: true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
								},
							},
						},
					},
				},
			},
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
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeObject,
								nullable: true,
								required: true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
										nullable: false,
										required: true,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "nested object slice with not required attrs",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{
							map[string]interface{}{"x": 1},
						},
					},
				},
				map[string]interface{}{
					"level1": map[string]interface{}{
						"level2": []interface{}{
							map[string]interface{}{"x": 1},
							map[string]interface{}{"x": 1, "y": "ok"},
							map[string]interface{}{"y": "thumbs up"},
						},
					},
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						nullable: false,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeObject,
								arrayLevel: 1,
								nullable:   false,
								required:   true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										nullable: false,
										required: false,
									},
									{
										key:      "y",
										t:        nodeTypeString,
										nullable: false,
										required: false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "nested object slice with not required attrs #2",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": true,
				},
				map[string]interface{}{
					"x": false,
					"y": []interface{}{
						map[string]interface{}{
							"a": "yes",
						},
						map[string]interface{}{
							"b": "no",
						},
					},
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeBool,
						nullable: false,
						required: true,
					},
					{
						key:        "y",
						t:          nodeTypeObject,
						arrayLevel: 1,
						nullable:   false,
						required:   false,
						children: []*node{
							{
								key:      "a",
								t:        nodeTypeString,
								nullable: false,
								required: false,
							},
							{
								key:      "b",
								t:        nodeTypeString,
								nullable: false,
								required: false,
							},
						},
					},
				},
			},
		},
		{
			name:        "list of objects with only null attr values",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": nil,
				},
				map[string]interface{}{
					"x": nil,
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: false,
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeInit, // only null values wont expand type
						nullable: true,
						required: true,
					},
				},
			},
		},
		{
			name:        "object + []object",
			startAsRoot: true,
			expands: []interface{}{
				map[string]interface{}{
					"x": true,
				},
				[]interface{}{
					map[string]interface{}{
						"x": true,
					},
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeInterface,
				nullable: false,
				required: true,
			},
		},
		{
			name:        "[]bool + bool",
			startAsRoot: true,
			expands: []interface{}{
				true,
				[]interface{}{true},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeInterface,
				arrayLevel: 0,
				nullable:   false,
				required:   true,
			},
		},
		{
			name:        "[][]bool + []bool",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true},
				[]interface{}{
					[]interface{}{true},
				},
			},
			expected: &node{
				root:       true,
				key:        baseTypeName,
				t:          nodeTypeInterface,
				arrayLevel: 0,
				nullable:   false,
				required:   true,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			f := newNode(baseTypeName)
			f.root = tc.startAsRoot

			for _, v := range tc.expands {
				f.grow(v)
			}
			f.sort()

			if !f.compare(tc.expected) {
				t.Errorf("invalid node\nwant:\n%s\ngot:\n%s", tc.expected.repr(""), f.repr(""))
			}
		})
	}
}

func TestJSONNodeTreeKeysNames(t *testing.T) {
	testCases := []struct {
		name          string
		growInput     interface{}
		expectedKeys  map[string]bool
		expectedNames map[string]bool
	}{
		{
			name: "",
			growInput: map[string]interface{}{
				"x_": 1,
				"x$": 2,
			},
			expectedKeys: map[string]bool{
				"x_": true,
				"x$": true,
			},
			expectedNames: map[string]bool{
				"X":  true,
				"X2": true,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			f := newNode(baseTypeName)
			f.root = true
			f.grow(tc.growInput)

			keys := make(map[string]bool)
			names := make(map[string]bool)
			for _, c := range f.children {
				keys[c.key] = true
				names[c.name] = true
			}
			assert.Equal(t, tc.expectedKeys, keys)
			assert.Equal(t, tc.expectedNames, names)
		})
	}
}

func TestJSONNodeExtractCommonSubtrees(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		root     *node
		expected []*node
	}{
		{
			name: "no extraction",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "fieldA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "fieldB",
						t:        nodeTypeObject,
						nullable: false,
						children: []*node{
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "fieldA",
							t:        nodeTypeObject,
							nullable: true,
							children: []*node{
								{
									key:      "x",
									t:        nodeTypeFloat,
									nullable: true,
								},
							},
						},
						{
							key:      "fieldB",
							t:        nodeTypeObject,
							nullable: false,
							children: []*node{
								{
									key:      "y",
									t:        nodeTypeFloat,
									nullable: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "extract one",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							nullable: true,
						},
					},
				},
			},
		},
		{
			name: "extract one not required, one required",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: false,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						nullable: false,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: false,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       false,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							nullable: false,
						},
					},
				},
			},
		},
		{
			name: "extract simple 2",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "size1",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "width",
								t:        nodeTypeInt,
								nullable: true,
							},
							{
								key:      "height",
								t:        nodeTypeInt,
								nullable: true,
							},
						},
					},
					{
						key:      "size2",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "width",
								t:        nodeTypeInt,
								nullable: true,
							},
							{
								key:      "height",
								t:        nodeTypeInt,
								nullable: true,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "size1",
							t:              nodeTypeExtracted,
							externalTypeID: "Size",
							nullable:       true,
						},
						{
							key:            "size2",
							t:              nodeTypeExtracted,
							externalTypeID: "Size",
							nullable:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							nullable: true,
						},
					},
				},
				{
					root:     true,
					key:      "size",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "height",
							t:        nodeTypeInt,
							nullable: true,
						},
						{
							key:      "width",
							t:        nodeTypeInt,
							nullable: true,
						},
					},
				},
			},
		},
		{
			name: "extract 2 with colliding names",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:      "pointC",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeInt,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeInt,
								nullable: true,
							},
						},
					},
					{
						key:      "pointD",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeInt,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeInt,
								nullable: true,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointC",
							t:              nodeTypeExtracted,
							externalTypeID: "Point2",
							nullable:       true,
						},
						{
							key:            "pointD",
							t:              nodeTypeExtracted,
							externalTypeID: "Point2",
							nullable:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							nullable: true,
						},
					},
				},
				{
					root:     true,
					key:      "point2",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeInt,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeInt,
							nullable: true,
						},
					},
				},
			},
		},
		{
			name: "extract object and slice of objects",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				nullable: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						nullable: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
					{
						key:        "pointsOther",
						t:          nodeTypeObject,
						nullable:   true,
						arrayLevel: 1,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								nullable: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								nullable: true,
							},
						},
					},
				},
			},
			expected: []*node{
				{
					root:     true,
					key:      baseTypeName,
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
						},
						{
							key:            "pointsOther",
							t:              nodeTypeExtracted,
							externalTypeID: "Point",
							nullable:       true,
							arrayLevel:     1,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					nullable: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							nullable: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							nullable: true,
						},
					},
				},
			},
		},
	}

	var setupNames func(n *node)
	setupNames = func(n *node) {
		n.name = attrName(n.key)
		for _, c := range n.children {
			setupNames(c)
		}
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// setup names here to skip making them by hand in test cases
			setupNames(tc.root)
			for _, n := range tc.expected {
				setupNames(n)
			}

			tc.root.sort()

			opts := options{}

			nodes := extractCommonSubtrees(tc.root)
			if !assert.Equal(t, len(tc.expected), len(nodes)) {
				t.Logf("\n%s\n\n", astPrintDecls(astMakeDecls(nodes, opts)))
				t.FailNow()
			}

			ok := true
			for i, n := range nodes {
				if !n.compare(tc.expected[i]) {
					t.Errorf("invald node %d, want:\n\n%s\n\ngot:\n\n%s", i, tc.expected[i].repr(""), n.repr(""))
					ok = false
				}
			}

			if !ok {
				t.Logf("\n%s\n\n", astPrintDecls(astMakeDecls(nodes, opts)))
			}
		})
	}
}

func TestArrayStructureDepth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		in               []interface{}
		inNodeType       nodeType
		expectedDepth    int
		expectedType     nodeType
		expectedNullable bool
	}{
		{
			name:          "empty array",
			in:            []interface{}{},
			expectedDepth: 1,
			expectedType:  nodeTypeInit,
		},
		{
			name:          "flat array",
			in:            []interface{}{1, 2, 3},
			expectedDepth: 1,
			expectedType:  nodeTypeInt,
		},
		{
			name:          "flat array, variable types",
			in:            []interface{}{1, "2", true},
			expectedDepth: 1,
			expectedType:  nodeTypeInterface,
		},
		{
			name:          "subarrays, variable structure",
			in:            []interface{}{1, 2, []interface{}{3, 4}},
			expectedDepth: -1,
			expectedType:  nodeTypeInterface,
		},
		{
			name: "subarrays, variable structure #2",
			in: []interface{}{
				[]interface{}{true, false},
				[]interface{}{
					[]interface{}{true, false},
					[]interface{}{true, false},
				},
			},
			expectedDepth: -1,
			expectedType:  nodeTypeInterface,
		},
		{
			name:          "subarrays, variable structure, variable types",
			in:            []interface{}{1, true, []interface{}{3, "4"}},
			expectedDepth: -1,
			expectedType:  nodeTypeInterface,
		},
		{
			name:          "subarrays, same structure, same types",
			in:            []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}},
			expectedDepth: 2,
			expectedType:  nodeTypeInt,
		},
		{
			name:          "subarrays, same structure, variable types",
			in:            []interface{}{[]interface{}{1, "2"}, []interface{}{3, true}},
			expectedDepth: 2,
			expectedType:  nodeTypeInterface,
		},
		{
			name:          "flat array, different in type",
			in:            []interface{}{1, 2, 3},
			inNodeType:    nodeTypeBool,
			expectedDepth: 1,
			expectedType:  nodeTypeInterface,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			d, tp, nullable := arrayStructure(tc.in, tc.inNodeType)
			assert.Equal(t, tc.expectedDepth, d)
			assert.Equal(t, tc.expectedType, tp)
			assert.Equal(t, tc.expectedNullable, nullable)
		})
	}
}

func TestArrayStructureType(t *testing.T) {
	tests := []struct {
		name         string
		in           []interface{}
		inType       nodeType
		wantType     nodeType
		wantNullable bool
	}{
		{
			name:         "no data",
			in:           []interface{}{},
			inType:       nil,
			wantType:     nodeTypeInit,
			wantNullable: false,
		},
		{
			name:         "empty array, type set already",
			in:           []interface{}{},
			inType:       nodeTypeString,
			wantType:     nodeTypeString,
			wantNullable: false,
		},
		{
			name:         "not empty array, type matches",
			in:           []interface{}{"a", "b"},
			inType:       nodeTypeString,
			wantType:     nodeTypeString,
			wantNullable: false,
		},
		{
			name: "not empty array, type matches, depth 2",
			in: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"a", "b"},
			},
			inType:       nodeTypeString,
			wantType:     nodeTypeString,
			wantNullable: false,
		},
		{
			name:         "not empty array, type not match",
			in:           []interface{}{"a", "b"},
			inType:       nodeTypeInt,
			wantType:     nodeTypeInterface,
			wantNullable: false,
		},
		{
			name: "not empty array, type not match, depth 2",
			in: []interface{}{
				[]interface{}{"a", "b"},
				[]interface{}{"a", "b"},
			},
			inType:       nodeTypeInt,
			wantType:     nodeTypeInterface,
			wantNullable: false,
		},
		{
			name:         "ints + floats",
			in:           []interface{}{1, 2, 3.14},
			inType:       nil,
			wantType:     nodeTypeFloat,
			wantNullable: false,
		},
		{
			name: "ints + floats, depth 2",
			in: []interface{}{
				[]interface{}{1, 2, 3.14},
				[]interface{}{1, 2, 3.14},
			},
			inType:       nil,
			wantType:     nodeTypeFloat,
			wantNullable: false,
		},
		{
			name:         "ints + null",
			in:           []interface{}{1, 2, 3, nil},
			inType:       nil,
			wantType:     nodeTypeInt,
			wantNullable: true,
		},
		{
			name:         "null + ints",
			in:           []interface{}{nil, 1, 2, 3},
			inType:       nil,
			wantType:     nodeTypeInt,
			wantNullable: true,
		},
		{
			name:         "nums + null",
			in:           []interface{}{1.45, nil, 1, 2.2, 3},
			inType:       nil,
			wantType:     nodeTypeFloat,
			wantNullable: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotType, nullable := arrayStructure(tt.in, tt.inType)
			assert.Equal(t, tt.wantType.id(), gotType.id())
			assert.Equal(t, tt.wantNullable, nullable)
		})
	}
}

func TestJSONNodeFits(t *testing.T) {
	objectV1_1 := newObjectNode(
		"test1",
		newObjectNode("cv11"),
		newObjectNode("cv12"),
	)
	objectV1_2 := newObjectNode(
		"test2",
		newObjectNode("cv11"),
		newObjectNode("cv12"),
		newObjectNode("cv13"),
	)
	objectV2_1 := newObjectNode(
		"test3",
		newObjectNode("cv21"),
		newObjectNode("cv22"),
		newObjectNode("cv23"),
	)
	objectV3_1 := newObjectNode(
		"test4",
		objectV1_1,
		objectV2_1,
		&node{t: nodeTypeBool, name: "cv31"},
		&node{t: nodeTypeFloat, name: "cv32"},
	)
	objectV3_2 := newObjectNode(
		"test5",
		objectV1_1,
		objectV1_2,
		objectV2_1,
		&node{t: nodeTypeBool, name: "cv31"},
		&node{t: nodeTypeFloat, name: "cv32"},
		&node{t: nodeTypeInt, name: "cv33"},
	)

	tests := []struct {
		name      string
		n1        *node
		n2        *node
		want      bool
		wantScore float64
	}{
		{
			name: "non object v1",
			n1:   &node{t: nodeTypeBool},
			n2:   objectV1_1,
			want: false,
		},
		{
			name: "non object v2",
			n1:   objectV1_1,
			n2:   &node{t: nodeTypeBool},
			want: false,
		},
		{
			name:      "v1_2 <= v1_1",
			n1:        objectV1_2,
			n2:        objectV1_1,
			want:      true,
			wantScore: 0.667, // v1_2 has 3 attributes, v1_1 has misses 1 attribute, (3-1)/3 == 0.667
		},
		{
			name: "v1_1 !<= v1_2",
			n1:   objectV1_1,
			n2:   objectV1_2,
			want: false,
		},
		{
			name: "v1_1 !<= v1_2",
			n1:   objectV1_1,
			n2:   objectV1_2,
			want: false,
		},
		{
			name: "v1_1 !<= v2_1",
			n1:   objectV1_1,
			n2:   objectV2_1,
			want: false,
		},
		{
			name: "v1_1 !=> v2_1",
			n1:   objectV2_1,
			n2:   objectV1_1,
			want: false,
		},
		{
			name:      "v3_1 => v3_2",
			n1:        objectV3_2,
			n2:        objectV3_1,
			want:      true,
			wantScore: 0.857, // v3_2 has 14 attributes, v3_1 misses 2 attributes, (14-2)/14 == 0.857
		},
		{
			name: "v3_1 !<= v3_2",
			n1:   objectV3_1,
			n2:   objectV3_2,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fits, score := tt.n1.fits(tt.n2)
			assert.Equal(t, tt.want, fits)
			assert.InDelta(t, tt.wantScore, score, 0.01)
			if fits {
				t.Log(score)
			}
		})
	}
}

func newObjectNode(name string, children ...*node) *node {
	return &node{
		root:           false,
		nullable:       false,
		required:       false,
		key:            name,
		name:           name,
		t:              nodeTypeObject,
		children:       children,
		arrayLevel:     0,
		arrayWithNulls: false,
	}
}
