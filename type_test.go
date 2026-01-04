package json2go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeExpand(t *testing.T) {
	testCases := []struct {
		name         string
		inputs       []any
		resultTypeID string
	}{
		// base types
		{
			name:         "input to bool",
			inputs:       []any{true},
			resultTypeID: nodeTypeBool.id(),
		},
		{
			name:         "input to int",
			inputs:       []any{1},
			resultTypeID: nodeTypeInt.id(),
		},
		{
			name:         "input to float",
			inputs:       []any{1.1},
			resultTypeID: nodeTypeFloat.id(),
		},
		{
			name:         "input to string",
			inputs:       []any{"123"},
			resultTypeID: nodeTypeString.id(),
		},
		{
			name:         "input to time",
			inputs:       []any{"2006-01-02T15:04:05+07:00"},
			resultTypeID: nodeTypeTime.id(),
		},
		{
			name: "input to object",
			inputs: []any{
				map[string]any{
					"key": "value",
				},
			},
			resultTypeID: nodeTypeObject.id(),
		},

		// mixed types
		{
			name:         "input to interface #1",
			inputs:       []any{"123", 123},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #2",
			inputs:       []any{123, "123"},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #3",
			inputs:       []any{"123", 123.4},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #4",
			inputs:       []any{123, true},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "input to interface #5",
			inputs:       []any{true, 123},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface #6",
			inputs: []any{
				map[string]any{
					"k": 1,
				},
				true,
				123,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - bool + []bool",
			inputs: []any{
				true,
				[]any{true},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []bool + bool",
			inputs: []any{
				[]any{true},
				true,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - int + []int",
			inputs: []any{
				1,
				[]any{1},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []int + int",
			inputs: []any{
				[]any{1},
				1,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - float + []float",
			inputs: []any{
				1.1,
				[]any{1.1},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []float + float",
			inputs: []any{
				[]any{1.1},
				1.1,
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - string + []string",
			inputs: []any{
				"1.1",
				[]any{"1.1"},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []string + string",
			inputs: []any{
				[]any{"1.1"},
				"1.1",
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name: "input to interface - []object + object",
			inputs: []any{
				[]any{
					map[string]any{
						"x": 1,
					},
				},
				map[string]any{
					"x": 1,
				},
			},
			resultTypeID: nodeTypeInterface.id(),
		},
		{
			name:         "time + string",
			inputs:       []any{"2006-01-02T15:04:05+07:00", "some stirng"},
			resultTypeID: nodeTypeString.id(),
		},
		{
			name:         "string + time",
			inputs:       []any{"some stirng", "2006-01-02T15:04:05+07:00"},
			resultTypeID: nodeTypeString.id(),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var k nodeType = nodeTypeInit
			for _, in := range tc.inputs {
				k = growType(k, in)
			}

			assert.Equal(t, tc.resultTypeID, k.id())
		})
	}
}
