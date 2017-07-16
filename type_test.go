package json2go

import "testing"

func TestTypeExpand(t *testing.T) {
	testCases := []struct {
		name         string
		inputs       []interface{}
		resultTypeID string
	}{
		// base types
		{
			name:         "input to bool",
			inputs:       []interface{}{true},
			resultTypeID: nodeTypeBool.id(),
		},
		{
			name:         "input to int",
			inputs:       []interface{}{1},
			resultTypeID: nodeTypeInt.id(),
		},
		{
			name:         "input to float",
			inputs:       []interface{}{1.1},
			resultTypeID: nodeTypeFloat.id(),
		},
		{
			name:         "input to string",
			inputs:       []interface{}{"123"},
			resultTypeID: nodeTypeString.id(),
		},
		{
			name: "input to object",
			inputs: []interface{}{
				map[string]interface{}{
					"key": "value",
				},
			},
			resultTypeID: nodeTypeObject.id(),
		},

		// mixed types
		{
			name:         "input to interface #1",
			inputs:       []interface{}{"123", 123},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #2",
			inputs:       []interface{}{123, "123"},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #3",
			inputs:       []interface{}{"123", 123.4},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #4",
			inputs:       []interface{}{123, true},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #5",
			inputs:       []interface{}{true, 123},
			resultTypeID: nodeTypeInterface.id(),
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
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - bool + []bool",
			inputs: []interface{}{
				true,
				[]interface{}{true},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []bool + bool",
			inputs: []interface{}{
				[]interface{}{true},
				true,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - int + []int",
			inputs: []interface{}{
				1,
				[]interface{}{1},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []int + int",
			inputs: []interface{}{
				[]interface{}{1},
				1,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - float + []float",
			inputs: []interface{}{
				1.1,
				[]interface{}{1.1},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []float + float",
			inputs: []interface{}{
				[]interface{}{1.1},
				1.1,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - string + []string",
			inputs: []interface{}{
				"1.1",
				[]interface{}{"1.1"},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []string + string",
			inputs: []interface{}{
				[]interface{}{"1.1"},
				"1.1",
			},
			resultTypeID: nodeTypeInterface.id(),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var k nodeType = nodeTypeInit
			for _, in := range tc.inputs {
				k = growType(k, in)
			}
			if k.id() != tc.resultTypeID {
				t.Errorf("invalid result type, want %s, got %s", tc.resultTypeID, k.id())
			}
		})
	}
}
