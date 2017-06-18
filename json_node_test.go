package json2go

import "testing"

func TestCompare(t *testing.T) {
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
				t:   newBoolType(),
			},
			n2: &node{
				key: "n2",
				t:   newBoolType(),
			},
			expectedEqual: false,
		},
		{
			name: "type not equal",
			n1: &node{
				key: "n",
				t:   newBoolType(),
			},
			n2: &node{
				key: "n",
				t:   newIntType(),
			},
			expectedEqual: false,
		},
		{
			name: "required not equal",
			n1: &node{
				key:      "n",
				t:        newBoolType(),
				required: true,
			},
			n2: &node{
				key:      "n",
				t:        newBoolType(),
				required: false,
			},
			expectedEqual: false,
		},
		{
			name: "simple equal",
			n1: &node{
				key: "n",
				t:   newBoolType(),
			},
			n2: &node{
				key: "n",
				t:   newBoolType(),
			},
			expectedEqual: true,
		},
		{
			name: "complex equal",
			n1: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
				},
			},
			n2: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
				},
			},
			expectedEqual: true,
		},
		{
			name: "complex - child num not equal",
			n1: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
				},
			},
			n2: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
				},
			},
			expectedEqual: false,
		},
		{
			name: "complex - child type not equal",
			n1: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newFloatArrayType(),
					},
				},
			},
			n2: &node{
				key: "n",
				t:   newBoolType(),
				children: []*node{
					{
						key: "n1",
						t:   newInterfaceArrayType(),
					},
				},
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

func TestNodeRepr(t *testing.T) {
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
				t:        newBoolType(),
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
				t:        newIntType(),
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
				t:        newFloatType(),
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
				t:        newStringType(),
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
				t:        newInterfaceType(),
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
				root:     true,
				key:      baseTypeName,
				t:        newBoolArrayType(),
				required: true,
			},
		},
		{
			name:        "[]int",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newIntArrayType(),
				required: true,
			},
		},
		{
			name:        "[]float",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1.1},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newFloatArrayType(),
				required: true,
			},
		},
		{
			name:        "[]string",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{"1.1"},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newStringArrayType(),
				required: true,
			},
		},
		{
			name:        "[]interface",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true, 1},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newInterfaceArrayType(),
				required: true,
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        newIntType(),
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
					"b": 2.0,
					"a": "str",
				},
				map[string]interface{}{
					"c": 12.7,
					"b": true,
					"a": "str2",
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "a",
						t:        newStringType(),
						required: true,
					},
					{
						key:      "b",
						t:        newInterfaceType(),
						required: true,
					},
					{
						key:      "c",
						t:        newFloatType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "slice",
						t:        newIntArrayType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newFloatArrayType(),
								required: true,
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        newIntType(),
						required: false,
					},
					{
						key:      "y",
						t:        newStringType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newObjectArrayType(),
								required: true,
								children: []*node{
									{
										key:      "x",
										t:        newIntType(),
										required: true,
									},
									{
										key:      "y",
										t:        newIntType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newUnknownObjectType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newObjectType(),
								required: false,
								children: []*node{
									{
										key:      "x",
										t:        newIntType(),
										required: true,
									},
									{
										key:      "y",
										t:        newIntType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newObjectType(),
								required: false,
								children: []*node{
									{
										key:      "x",
										t:        newIntType(),
										required: true,
									},
									{
										key:      "y",
										t:        newIntType(),
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
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "level1",
						t:        newObjectType(),
						required: true,
						children: []*node{
							{
								key:      "level2",
								t:        newObjectArrayType(),
								required: true,
								children: []*node{
									{
										key:      "x",
										t:        newIntType(),
										required: false,
									},
									{
										key:      "y",
										t:        newStringType(),
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
						map[string]interface{}{"a": "yes"},
						map[string]interface{}{"b": "no"},
					},
				},
			},
			expected: &node{
				root:     true,
				key:      baseTypeName,
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        newBoolType(),
						required: true,
					},
					{
						key:      "y",
						t:        newObjectArrayType(),
						required: false,
						children: []*node{
							{
								key:      "a",
								t:        newStringType(),
								required: false,
							},
							{
								key:      "b",
								t:        newStringType(),
								required: false,
							},
						},
					},
				},
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
				t.Fatalf("invalid node. want:\n%s\ngot:\n%s", tc.expected.repr(""), f.repr(""))
			}
		})
	}
}

func TestTreeInfo(t *testing.T) {
	n := &node{
		root:     true,
		key:      baseTypeName,
		t:        newObjectType(),
		required: true,
		children: []*node{
			{
				key:      "pointA",
				t:        newObjectType(),
				required: true,
				children: []*node{
					{
						key:      "x",
						t:        newFloatType(),
						required: true,
					},
					{
						key:      "y",
						t:        newFloatType(),
						required: true,
					},
				},
			},
			{
				key:      "pointB",
				t:        newObjectType(),
				required: false,
				children: []*node{
					{
						key:      "x",
						t:        newFloatType(),
						required: true,
					},
					{
						key:      "y",
						t:        newFloatType(),
						required: true,
					},
				},
			},
		},
	}

	infos := make(map[string]nodeStructureInfo)
	n.treeInfo(infos)

	for k, info := range infos {
		t.Logf("%s: %s - %v\n", k, info.typeID, info.nodes)
		if len(info.nodes) > 1 && info.typeID == nodeTypeObject {
			extractedNode := *info.nodes[0]
			extractedNode.key = "point"

			n.modify(info.structureID, func(modNode *node) {
				modNode.t = newExternalObjectType()
				modNode.externalTypeID = "Point"
			})

			t.Logf("\n%s\n\n", astPrintDecls(astMakeDecls([]*node{n, &extractedNode})))
		}
	}

}
