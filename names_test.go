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
			name:         "long snake",
			fieldName:    "field__name",
			expectedName: "FieldName",
		},
		{
			name:         "snake at ends",
			fieldName:    "_field_name_",
			expectedName: "FieldName",
		},
		{
			name:         "long snake at ends",
			fieldName:    "__field_name",
			expectedName: "FieldName",
		},
		{
			name:         "snake case with initialism",
			fieldName:    "html_field_name",
			expectedName: "HTMLFieldName",
		},
		{
			name:         "mixed case",
			fieldName:    "aBcDeFg",
			expectedName: "Abcdefg",
		},
		{
			name:         "special chars",
			fieldName:    "$field_$name日本語",
			expectedName: "FieldName",
		},
		{
			name:         "garbage",
			fieldName:    "$@!%^&*()",
			expectedName: "",
		},
		{
			name:         "starting with digits",
			fieldName:    "123key",
			expectedName: "Key",
		},
		{
			name:         "name with digits",
			fieldName:    "key_666",
			expectedName: "Key666",
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
