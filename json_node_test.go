package json2go

import "testing"

func TestJSONNodeCompare(t *testing.T) {
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
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeInt,
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
				required: true,
				children: []*node{
					{
						key:      "a",
						t:        nodeTypeString,
						required: true,
					},
					{
						key:      "b",
						t:        nodeTypeInterface,
						required: true,
					},
					{
						key:      "c",
						t:        nodeTypeFloat,
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
				required: true,
				children: []*node{
					{
						key:        "slice",
						t:          nodeTypeInt,
						arrayLevel: 1,
						required:   true,
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeFloat,
								arrayLevel: 1,
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
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeInt,
						required: false,
					},
					{
						key:      "y",
						t:        nodeTypeString,
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeObject,
								arrayLevel: 1,
								required:   true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeUnknownObject,
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeObject,
								required: false,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        nodeTypeObject,
								required: false,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										required: true,
									},
									{
										key:      "y",
										t:        nodeTypeInt,
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
			name:        "nested object slice with nullable attrs",
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
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:        "level2",
								t:          nodeTypeObject,
								arrayLevel: 1,
								required:   true,
								children: []*node{
									{
										key:      "x",
										t:        nodeTypeInt,
										required: false,
									},
									{
										key:      "y",
										t:        nodeTypeString,
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
			name:        "nested object slice with nullable attrs #2",
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
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        nodeTypeBool,
						required: true,
					},
					{
						key:        "y",
						t:          nodeTypeObject,
						arrayLevel: 1,
						required:   false,
						children: []*node{
							{
								key:      "a",
								t:        nodeTypeString,
								required: false,
							},
							{
								key:      "b",
								t:        nodeTypeString,
								required: false,
							},
						},
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
				required:   true,
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

func TestJSONNodeExtractCommonSubtrees(t *testing.T) {
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
				required: true,
				children: []*node{
					{
						key:      "fieldA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "fieldB",
						t:        nodeTypeObject,
						required: false,
						children: []*node{
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
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
					required: true,
					children: []*node{
						{
							key:      "fieldA",
							t:        nodeTypeObject,
							required: true,
							children: []*node{
								{
									key:      "x",
									t:        nodeTypeFloat,
									required: true,
								},
							},
						},
						{
							key:      "fieldB",
							t:        nodeTypeObject,
							required: false,
							children: []*node{
								{
									key:      "y",
									t:        nodeTypeFloat,
									required: true,
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
				required: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
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
					required: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							required: true,
						},
					},
				},
			},
		},
		{
			name: "extract one required, one not required",
			root: &node{
				root:     true,
				key:      baseTypeName,
				t:        nodeTypeObject,
				required: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: false,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						required: false,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: false,
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
					required: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       false,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							required: false,
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
				required: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "size1",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "width",
								t:        nodeTypeInt,
								required: true,
							},
							{
								key:      "height",
								t:        nodeTypeInt,
								required: true,
							},
						},
					},
					{
						key:      "size2",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "width",
								t:        nodeTypeInt,
								required: true,
							},
							{
								key:      "height",
								t:        nodeTypeInt,
								required: true,
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
					required: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "size1",
							t:              nodeTypeExternal,
							externalTypeID: "Size",
							required:       true,
						},
						{
							key:            "size2",
							t:              nodeTypeExternal,
							externalTypeID: "Size",
							required:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							required: true,
						},
					},
				},
				{
					root:     true,
					key:      "size",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "height",
							t:        nodeTypeInt,
							required: true,
						},
						{
							key:      "width",
							t:        nodeTypeInt,
							required: true,
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
				required: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "pointB",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:      "pointC",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeInt,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeInt,
								required: true,
							},
						},
					},
					{
						key:      "pointD",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeInt,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeInt,
								required: true,
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
					required: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointB",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointC",
							t:              nodeTypeExternal,
							externalTypeID: "Point2",
							required:       true,
						},
						{
							key:            "pointD",
							t:              nodeTypeExternal,
							externalTypeID: "Point2",
							required:       true,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							required: true,
						},
					},
				},
				{
					root:     true,
					key:      "point2",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeInt,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeInt,
							required: true,
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
				required: true,
				children: []*node{
					{
						key:      "pointA",
						t:        nodeTypeObject,
						required: true,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
							},
						},
					},
					{
						key:        "pointsOther",
						t:          nodeTypeObject,
						required:   true,
						arrayLevel: 1,
						children: []*node{
							{
								key:      "x",
								t:        nodeTypeFloat,
								required: true,
							},
							{
								key:      "y",
								t:        nodeTypeFloat,
								required: true,
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
					required: true,
					children: []*node{
						{
							key:            "pointA",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
						},
						{
							key:            "pointsOther",
							t:              nodeTypeExternal,
							externalTypeID: "Point",
							required:       true,
							arrayLevel:     1,
						},
					},
				},
				{
					root:     true,
					key:      "point",
					t:        nodeTypeObject,
					required: true,
					children: []*node{
						{
							key:      "x",
							t:        nodeTypeFloat,
							required: true,
						},
						{
							key:      "y",
							t:        nodeTypeFloat,
							required: true,
						},
					},
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.root.sort()

			nodes := extractCommonSubtrees(tc.root)
			if len(nodes) != len(tc.expected) {
				t.Logf("\n%s\n\n", astPrintDecls(astMakeDecls(nodes)))
				t.Fatalf("got invalid num of nodes, want %d, got %d", len(tc.expected), len(nodes))
			}

			ok := true
			for i, n := range nodes {
				if !n.compare(tc.expected[i]) {
					t.Errorf("invald node %d, want:\n\n%s\n\ngot:\n\n%s", i, tc.expected[i].repr(""), n.repr(""))
					ok = false
				}
			}

			if !ok {
				t.Logf("\n%s\n\n", astPrintDecls(astMakeDecls(nodes)))
			}
		})
	}
}

func TestArrayDepth(t *testing.T) {
	testCases := []struct {
		name          string
		in            []interface{}
		inNodeType    nodeType
		expectedDepth int
		expectedType  nodeType
	}{
		{
			name:          "empty array",
			in:            []interface{}{},
			expectedDepth: 1,
			expectedType:  nodeTypeInterface,
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
			d, tp := arrayStructure(tc.in, tc.inNodeType)
			if d != tc.expectedDepth {
				t.Errorf("want: %d, got %d", tc.expectedDepth, d)
			}
			if tp != tc.expectedType {
				t.Errorf("want: %v, got %v", tc.expectedType, tp)
			}
		})
	}
}
