package json2go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNames(t *testing.T) {
	t.Parallel()

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
			name:         "camel case",
			fieldName:    "camelCaseName",
			expectedName: "CamelCaseName",
		},
		{
			name:         "mixed case",
			fieldName:    "mixed_caseName_test",
			expectedName: "MixedCaseNameTest",
		},
		{
			name:         "mixed case with initialism",
			fieldName:    "mixed_caseName_htmlTest",
			expectedName: "MixedCaseNameHTMLTest",
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
			assert.Equal(t, tc.expectedName, attrName(tc.fieldName))
		})
	}
}

func TestExtractCommonName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		inputs   []string
		expected string
	}{
		{
			name:     "empty",
			inputs:   []string{},
			expected: "",
		},
		{
			name:     "single input",
			inputs:   []string{"test"},
			expected: "test",
		},
		{
			name:     "multiple same strings",
			inputs:   []string{"test", "test", "test"},
			expected: "test",
		},
		{
			name:     "common prefix",
			inputs:   []string{"test1", "test2", "testXyz", "test"},
			expected: "test",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, extractCommonName(tc.inputs...))
		})
	}
}

func TestNameFromNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		inputs   []string
		expected string
	}{
		{
			name:     "empty",
			inputs:   []string{},
			expected: "",
		},
		{
			name:     "single input",
			inputs:   []string{"test"},
			expected: "test",
		},
		{
			name:     "multiple long strings",
			inputs:   []string{"aaaaaa", "bbbbbb", "cccccc", "dddddd"},
			expected: "aaaaaa_bbbbbb_cccccc",
		},
		{
			name:     "multiple short strings",
			inputs:   []string{"a", "b", "c", "d", "e", "f"},
			expected: "a_b_c",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, nameFromNames(tc.inputs...))
		})
	}
}

func TestNextName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "1",
		},
		{
			name:     "simple",
			input:    "simple",
			expected: "simple2",
		},
		{
			name:     "simple2",
			input:    "simple2",
			expected: "simple3",
		},
		{
			name:     "numbers in name",
			input:    "2sim234ple3",
			expected: "2sim234ple4",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, nextName(tc.input))
		})
	}
}
