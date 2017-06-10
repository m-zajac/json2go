package jsontogo

import "testing"

func TestKindExpand(t *testing.T) {
	testCases := []struct {
		name           string
		inputs         []interface{}
		resultKindName kindName
	}{
		// base kinds
		{
			name:           "input to bool",
			inputs:         []interface{}{true},
			resultKindName: kindNameBool,
		},
		{
			name:           "input to int",
			inputs:         []interface{}{1},
			resultKindName: kindNameInt,
		},
		{
			name:           "input to float",
			inputs:         []interface{}{1.1},
			resultKindName: kindNameFloat,
		},
		{
			name:           "input to string",
			inputs:         []interface{}{"123"},
			resultKindName: kindNameString,
		},

		// mixed types
		{
			name:           "input to interface #1",
			inputs:         []interface{}{"123", 123},
			resultKindName: kindNameInterface,
		},
		{
			name:           "input to interface #2",
			inputs:         []interface{}{123, "123"},
			resultKindName: kindNameInterface,
		},
		{
			name:           "input to interface #3",
			inputs:         []interface{}{"123", 123.4},
			resultKindName: kindNameInterface,
		},
		{
			name:           "input to interface #4",
			inputs:         []interface{}{123, true},
			resultKindName: kindNameInterface,
		},
		{
			name:           "input to interface #5",
			inputs:         []interface{}{true, 123},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - bool + []bool",
			inputs: []interface{}{
				true,
				[]interface{}{true},
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - []bool + bool",
			inputs: []interface{}{
				[]interface{}{true},
				true,
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - int + []int",
			inputs: []interface{}{
				1,
				[]interface{}{1},
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - []int + int",
			inputs: []interface{}{
				[]interface{}{1},
				1,
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - float + []float",
			inputs: []interface{}{
				1.1,
				[]interface{}{1.1},
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - []float + float",
			inputs: []interface{}{
				[]interface{}{1.1},
				1.1,
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - string + []string",
			inputs: []interface{}{
				"1.1",
				[]interface{}{"1.1"},
			},
			resultKindName: kindNameInterface,
		},
		{
			name: "input to interface - []string + string",
			inputs: []interface{}{
				[]interface{}{"1.1"},
				"1.1",
			},
			resultKindName: kindNameInterface,
		},

		// arrays
		{
			name: "input to []bool",
			inputs: []interface{}{
				[]interface{}{true, false},
			},
			resultKindName: kindNameArrayBool,
		},
		{
			name: "input to []int",
			inputs: []interface{}{
				[]interface{}{1, 2},
			},
			resultKindName: kindNameArrayInt,
		},
		{
			name: "input to []float",
			inputs: []interface{}{
				[]interface{}{1.1, 2.2},
			},
			resultKindName: kindNameArrayFloat,
		},
		{
			name: "input to []float #2",
			inputs: []interface{}{
				[]interface{}{1, 2.2},
			},
			resultKindName: kindNameArrayFloat,
		},
		{
			name: "input to []float #3",
			inputs: []interface{}{
				[]interface{}{1.1, 2},
			},
			resultKindName: kindNameArrayFloat,
		},
		{
			name: "input to []string",
			inputs: []interface{}{
				[]interface{}{"xxx", "yyy"},
			},
			resultKindName: kindNameArrayString,
		},
		{
			name: "input to []interface{}",
			inputs: []interface{}{
				[]interface{}{true, 1},
			},
			resultKindName: kindNameArrayInterface,
		},
		{
			name: "input to []interface{} #2",
			inputs: []interface{}{
				[]interface{}{true, 1},
				[]interface{}{false, 2},
			},
			resultKindName: kindNameArrayInterface,
		},
		{
			name: "input to []interface{} #3",
			inputs: []interface{}{
				[]interface{}{true, 1},
				[]interface{}{2, false},
			},
			resultKindName: kindNameArrayInterface,
		},
		{
			name: "input to {}interface{} - []bool + []int",
			inputs: []interface{}{
				[]interface{}{true, false},
				[]interface{}{1, 2},
			},
			resultKindName: kindNameArrayInterface,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			k := newStartingKind()
			for _, in := range tc.inputs {
				k = k.grow(in)
			}
			if k.name != tc.resultKindName {
				t.Errorf("invalid bool kind expand type for bool input, want %s, got %s", tc.resultKindName, k.name)
			}
		})
	}
}
