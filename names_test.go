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

func TestNameFromNamesCapped(t *testing.T) {
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
			name:     "exactly max parts",
			inputs:   []string{"a", "b", "c"},
			expected: "a_b_c",
		},
		{
			name:     "more than max parts (4)",
			inputs:   []string{"a", "b", "c", "d"},
			expected: "a_b_cAnd1More",
		},
		{
			name:     "more than max parts (5)",
			inputs:   []string{"a", "b", "c", "d", "e"},
			expected: "a_b_cAnd2More",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, nameFromNamesCapped(tc.inputs...))
		})
	}
}

func TestSingularize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
		{
			name:     "regular plural - details",
			input:    "Details",
			expected: "Detail",
		},
		{
			name:     "regular plural - items",
			input:    "Items",
			expected: "Item",
		},
		{
			name:     "regular plural - users",
			input:    "Users",
			expected: "User",
		},
		{
			name:     "regular plural - products",
			input:    "Products",
			expected: "Product",
		},
		{
			name:     "plural -ies - categories",
			input:    "Categories",
			expected: "Category",
		},
		{
			name:     "plural -ies - entries",
			input:    "Entries",
			expected: "Entry",
		},
		{
			name:     "plural -sses - addresses",
			input:    "Addresses",
			expected: "Address",
		},
		{
			name:     "plural -xes - boxes",
			input:    "Boxes",
			expected: "Box",
		},
		{
			name:     "plural -ches - matches",
			input:    "Matches",
			expected: "Match",
		},
		{
			name:     "plural -shes - wishes",
			input:    "Wishes",
			expected: "Wish",
		},
		{
			name:     "too short - as",
			input:    "As",
			expected: "As",
		},
		{
			name:     "too short - is",
			input:    "Is",
			expected: "Is",
		},
		{
			name:     "not plural - item",
			input:    "Item",
			expected: "Item",
		},
		{
			name:     "irregular - children",
			input:    "Children",
			expected: "Child",
		},
		{
			name:     "irregular - people",
			input:    "People",
			expected: "Person",
		},
		{
			name:     "irregular - men",
			input:    "Men",
			expected: "Man",
		},
		{
			name:     "irregular - women",
			input:    "Women",
			expected: "Woman",
		},
		{
			name:     "irregular - teeth",
			input:    "Teeth",
			expected: "Tooth",
		},
		{
			name:     "irregular - mice",
			input:    "Mice",
			expected: "Mouse",
		},
		{
			name:     "irregular - indices",
			input:    "Indices",
			expected: "Index",
		},
		{
			name:     "uncountable - sheep",
			input:    "Sheep",
			expected: "Sheep",
		},
		{
			name:     "uncountable - deer",
			input:    "Deer",
			expected: "Deer",
		},
		{
			name:     "uncountable - fish",
			input:    "Fish",
			expected: "Fish",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, singularize(tc.input))
		})
	}
}

func TestTypeNameFromFieldName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		fieldName string
		expected  string
	}{
		{
			name:      "singular field",
			fieldName: "user",
			expected:  "User",
		},
		{
			name:      "plural field - details",
			fieldName: "details",
			expected:  "Detail",
		},
		{
			name:      "plural field - items",
			fieldName: "items",
			expected:  "Item",
		},
		{
			name:      "snake case plural - user_details",
			fieldName: "user_details",
			expected:  "UserDetail",
		},
		{
			name:      "irregular plural - children",
			fieldName: "children",
			expected:  "Child",
		},
		{
			name:      "irregular plural - people",
			fieldName: "people",
			expected:  "Person",
		},
		{
			name:      "camel case plural",
			fieldName: "userItems",
			expected:  "UserItem",
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, typeNameFromFieldName(tc.fieldName))
		})
	}
}
