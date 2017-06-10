package json2go

import "testing"

func TestTypeExpand(t *testing.T) {
	testCases := []struct {
		name         string
		inputs       []interface{}
		resultTypeID fieldTypeID
	}{
		// base types
		{
			name:         "input to bool",
			inputs:       []interface{}{true},
			resultTypeID: fieldTypeBool,
		},
		{
			name:         "input to int",
			inputs:       []interface{}{1},
			resultTypeID: fieldTypeInt,
		},
		{
			name:         "input to float",
			inputs:       []interface{}{1.1},
			resultTypeID: fieldTypeFloat,
		},
		{
			name:         "input to string",
			inputs:       []interface{}{"123"},
			resultTypeID: fieldTypeString,
		},
		{
			name: "input to object",
			inputs: []interface{}{
				map[string]interface{}{
					"key": "value",
				},
			},
			resultTypeID: fieldTypeObject,
		},

		// mixed types
		{
			name:         "input to interface #1",
			inputs:       []interface{}{"123", 123},
			resultTypeID: fieldTypeInterface,
		},
		{
			name:         "input to interface #2",
			inputs:       []interface{}{123, "123"},
			resultTypeID: fieldTypeInterface,
		},
		{
			name:         "input to interface #3",
			inputs:       []interface{}{"123", 123.4},
			resultTypeID: fieldTypeInterface,
		},
		{
			name:         "input to interface #4",
			inputs:       []interface{}{123, true},
			resultTypeID: fieldTypeInterface,
		},
		{
			name:         "input to interface #5",
			inputs:       []interface{}{true, 123},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface #6",
			inputs: []interface{}{
				map[string]interface{}{
					"k": 1,
				},
				true,
				123,
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - bool + []bool",
			inputs: []interface{}{
				true,
				[]interface{}{true},
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - []bool + bool",
			inputs: []interface{}{
				[]interface{}{true},
				true,
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - int + []int",
			inputs: []interface{}{
				1,
				[]interface{}{1},
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - []int + int",
			inputs: []interface{}{
				[]interface{}{1},
				1,
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - float + []float",
			inputs: []interface{}{
				1.1,
				[]interface{}{1.1},
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - []float + float",
			inputs: []interface{}{
				[]interface{}{1.1},
				1.1,
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - string + []string",
			inputs: []interface{}{
				"1.1",
				[]interface{}{"1.1"},
			},
			resultTypeID: fieldTypeInterface,
		},
		{
			name: "input to interface - []string + string",
			inputs: []interface{}{
				[]interface{}{"1.1"},
				"1.1",
			},
			resultTypeID: fieldTypeInterface,
		},

		// arrays
		{
			name: "input to []bool",
			inputs: []interface{}{
				[]interface{}{true, false},
			},
			resultTypeID: fieldTypeArrayBool,
		},
		{
			name: "input to []int",
			inputs: []interface{}{
				[]interface{}{1, 2},
				[]interface{}{4, 5, 6},
			},
			resultTypeID: fieldTypeArrayInt,
		},
		{
			name: "input to []int with float types",
			inputs: []interface{}{
				[]interface{}{1.0, 2.0},
				[]interface{}{4.0, 5.0, 6.0},
			},
			resultTypeID: fieldTypeArrayInt,
		},
		{
			name: "input to []float",
			inputs: []interface{}{
				[]interface{}{1.1, 2.2},
			},
			resultTypeID: fieldTypeArrayFloat,
		},
		{
			name: "input to []float #2",
			inputs: []interface{}{
				[]interface{}{1, 2.2},
			},
			resultTypeID: fieldTypeArrayFloat,
		},
		{
			name: "input to []float #3",
			inputs: []interface{}{
				[]interface{}{1.1, 2},
			},
			resultTypeID: fieldTypeArrayFloat,
		},
		{
			name: "input to []string",
			inputs: []interface{}{
				[]interface{}{"xxx", "yyy"},
			},
			resultTypeID: fieldTypeArrayString,
		},
		{
			name: "input to []interface{}",
			inputs: []interface{}{
				[]interface{}{true, 1},
			},
			resultTypeID: fieldTypeArrayInterface,
		},
		{
			name: "input to []interface{} #2",
			inputs: []interface{}{
				[]interface{}{true, 1},
				[]interface{}{false, 2},
			},
			resultTypeID: fieldTypeArrayInterface,
		},
		{
			name: "input to []interface{} #3",
			inputs: []interface{}{
				[]interface{}{true, 1},
				[]interface{}{2, false},
			},
			resultTypeID: fieldTypeArrayInterface,
		},
		{
			name: "input to {}interface{} - []bool + []int",
			inputs: []interface{}{
				[]interface{}{true, false},
				[]interface{}{1, 2},
			},
			resultTypeID: fieldTypeArrayInterface,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			k := newInitType()
			for _, in := range tc.inputs {
				k = k.grow(in)
			}
			if k.id != tc.resultTypeID {
				t.Errorf("invalid result type, want %s, got %s", tc.resultTypeID, k.id)
			}
		})
	}
}
