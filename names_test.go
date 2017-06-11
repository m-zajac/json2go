package json2go

import "testing"

func TestNames(t *testing.T) {
	testCases := []struct {
		name         string
		fieldName    string
		expectedName string
	}{
		{
			name:         "simple lowercase",
			fieldName:    "name",
			expectedName: "Name",
		},
		{
			name:         "snake case",
			fieldName:    "field_name",
			expectedName: "FieldName",
		},
		{
			name:         "snake case with initialism",
			fieldName:    "html_field_name",
			expectedName: "HTMLFieldName",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			if name := attrName(tc.fieldName); name != tc.expectedName {
				t.Errorf("invalid name, want `%s`, got `%s`", tc.expectedName, name)
			}
		})
	}
}
