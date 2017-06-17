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
				name: "n",
				t:    newBoolType(),
			},
			n2: &node{
				name: "n2",
				t:    newBoolType(),
			},
			expectedEqual: false,
		},
		{
			name: "type not equal",
			n1: &node{
				name: "n",
				t:    newBoolType(),
			},
			n2: &node{
				name: "n",
				t:    newIntType(),
			},
			expectedEqual: false,
		},
		{
			name: "simple equal",
			n1: &node{
				name: "n",
				t:    newBoolType(),
			},
			n2: &node{
				name: "n",
				t:    newBoolType(),
			},
			expectedEqual: true,
		},
		{
			name: "complex equal",
			n1: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newFloatArrayType(),
					},
				},
			},
			n2: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newFloatArrayType(),
					},
				},
			},
			expectedEqual: true,
		},
		{
			name: "complex - child num not equal",
			n1: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newFloatArrayType(),
					},
				},
			},
			n2: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newFloatArrayType(),
					},
					"n2": {
						name: "n1",
						t:    newFloatArrayType(),
					},
				},
			},
			expectedEqual: false,
		},
		{
			name: "complex - child type not equal",
			n1: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newFloatArrayType(),
					},
				},
			},
			n2: &node{
				name: "n",
				t:    newBoolType(),
				children: map[string]*node{
					"n1": {
						name: "n1",
						t:    newInterfaceArrayType(),
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
				root: true,
				name: baseTypeName,
				t:    newBoolType(),
			},
		},
		{
			name:        "int",
			startAsRoot: true,
			expands:     []interface{}{1},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newIntType(),
			},
		},
		{
			name:        "float",
			startAsRoot: true,
			expands:     []interface{}{1.1},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newFloatType(),
			},
		},
		{
			name:        "string",
			startAsRoot: true,
			expands:     []interface{}{"1.1"},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newStringType(),
			},
		},
		{
			name:        "interface",
			startAsRoot: true,
			expands:     []interface{}{true, 1},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newInterfaceType(),
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
				root: true,
				name: baseTypeName,
				t:    newBoolArrayType(),
			},
		},
		{
			name:        "[]int",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1},
			},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newIntArrayType(),
			},
		},
		{
			name:        "[]float",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{1.1},
			},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newFloatArrayType(),
			},
		},
		{
			name:        "[]string",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{"1.1"},
			},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newStringArrayType(),
			},
		},
		{
			name:        "[]interface",
			startAsRoot: true,
			expands: []interface{}{
				[]interface{}{true, 1},
			},
			expected: &node{
				root: true,
				name: baseTypeName,
				t:    newInterfaceArrayType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"x": {
						name: "x",
						t:    newIntType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"a": {
						name: "a",
						t:    newStringType(),
					},
					"b": {
						name: "b",
						t:    newInterfaceType(),
					},
					"c": {
						name: "c",
						t:    newFloatType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"slice": {
						name: "slice",
						t:    newIntArrayType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"level1": {
						name: "level1",
						t:    newObjectType(),
						children: map[string]*node{
							"level2": {
								name: "level2",
								t:    newFloatArrayType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"x": {
						name: "x",
						t:    newIntType(),
					},
					"y": {
						name: "y",
						t:    newStringType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"level1": {
						name: "level1",
						t:    newObjectType(),
						children: map[string]*node{
							"level2": {
								name: "level2",
								t:    newObjectArrayType(),
								children: map[string]*node{
									"x": {
										name: "x",
										t:    newIntType(),
									},
									"y": {
										name: "y",
										t:    newIntType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"level1": {
						name: "level1",
						t:    newObjectType(),
						children: map[string]*node{
							"level2": {
								name: "level2",
								t:    newUnknownObjectType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"level1": {
						name: "level1",
						t:    newObjectType(),
						children: map[string]*node{
							"level2": {
								name: "level2",
								t:    newObjectType(),
								children: map[string]*node{
									"x": {
										name: "x",
										t:    newIntType(),
									},
									"y": {
										name: "y",
										t:    newIntType(),
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
				root: true,
				name: baseTypeName,
				t:    newObjectType(),
				children: map[string]*node{
					"level1": {
						name: "level1",
						t:    newObjectType(),
						children: map[string]*node{
							"level2": {
								name: "level2",
								t:    newObjectType(),
								children: map[string]*node{
									"x": {
										name: "x",
										t:    newIntType(),
									},
									"y": {
										name: "y",
										t:    newIntType(),
									},
								},
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

			if !f.compare(tc.expected) {
				t.Fatalf("invalid node. want:\n%s\ngot:\n%s", tc.expected.repr(""), f.repr(""))
			}
		})
	}
}
