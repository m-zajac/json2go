package json2go

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"
)

const baseTypeName = "Object"

type field struct {
	root   bool
	name   string
	t      fieldType
	fields map[string]*field
}

func newField(name string) *field {
	return &field{
		name:   name,
		t:      newInitType(),
		fields: make(map[string]*field),
	}
}

func (f *field) grow(input interface{}) {
	f.t = f.t.grow(input)

	switch typedInput := input.(type) {
	case map[string]interface{}:
		for k, v := range typedInput {
			if _, ok := f.fields[k]; !ok {
				f.fields[k] = newField(k)
			}
			f.fields[k].grow(v)
		}
	case []interface{}:
		if f.t.id != fieldTypeArrayObject {
			break
		}

	loop:
		for _, iv := range typedInput {
			mv, ok := iv.(map[string]interface{})
			if !ok {
				break loop
			}

			for k, v := range mv {
				if _, ok := f.fields[k]; !ok {
					f.fields[k] = newField(k)
				}
				f.fields[k].grow(v)
			}

		}

	}
}

func (f *field) repr() (string, error) {
	var b bytes.Buffer

	// write valid go file header
	fileHeader := "package main\n\n"
	b.WriteString(fileHeader)

	// write struct representation
	f.reprToBuffer(&b)

	// parse and print formatted
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "", b.String(), 0)
	if err != nil {
		return b.String(), fmt.Errorf("invalid type def, parser error: %v", err)
	}

	b.Reset()
	printer.Fprint(&b, fset, astf)

	// remove go file header
	repr := b.String()
	repr = strings.TrimPrefix(repr, fileHeader)
	repr = strings.TrimSpace(repr)

	return repr, nil
}

func (f *field) reprToBuffer(b *bytes.Buffer) {
	if f.root {
		b.WriteString(fmt.Sprintf("type %s ", baseTypeName))
	} else {
		name := strings.Title(f.name)
		b.WriteString(fmt.Sprintf("%s ", name))
	}
	b.WriteString(f.t.repr())

	if len(f.fields) > 0 && (f.t.id == fieldTypeObject || f.t.id == fieldTypeArrayObject) {
		b.WriteString(" {")
		defer b.WriteString("}")

		// sort subfields by name
		type fieldWithName struct {
			name  string
			field *field
		}
		var sortedFields []fieldWithName
		for subName, subField := range f.fields {
			sortedFields = append(sortedFields, fieldWithName{
				name:  subName,
				field: subField,
			})
		}
		sort.Slice(sortedFields, func(i, j int) bool {
			return sortedFields[i].name < sortedFields[j].name
		})

		// repr subfields
		for _, subField := range sortedFields {
			b.WriteString("\n")

			// repr subfield
			subField.field.reprToBuffer(b)

			// add json tag
			b.WriteString(fmt.Sprintf(" `json:\"%s\"`", subField.name))
		}
		b.WriteString("\n")
	}
}
