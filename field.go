package jsontogo

import (
	"fmt"
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
	switch t := input.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if _, ok := f.fields[k]; !ok {
				f.fields[k] = newField(k)
			}
			f.fields[k].grow(v)
		}
	default:
		f.t = f.t.grow(input)
	}
}

func (f *field) repr() string {
	if len(f.fields) == 0 {
		if f.root {
			return fmt.Sprintf("type %s %s", baseTypeName, f.t.repr())
		}
		return fmt.Sprintf("%s %s", strings.Title(f.name), f.t.repr())
	}

	return ""
}
