package json2go

import (
	"bytes"
	"errors"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"sort"
	"strings"
)

const baseTypeName = "Object"

type node struct {
	root     bool
	name     string
	t        nodeType
	children map[string]*node
}

func newNode(name string) *node {
	return &node{
		name:     name,
		t:        newInitType(),
		children: make(map[string]*node),
	}
}

func (f *node) grow(input interface{}) {
	f.t = f.t.grow(input)

	switch typedInput := input.(type) {
	case map[string]interface{}:
		for k, v := range typedInput {
			if _, ok := f.children[k]; !ok {
				f.children[k] = newNode(k)
			}
			f.children[k].grow(v)
		}
	case []interface{}:
		if f.t.id != nodeTypeArrayObject {
			break
		}

	loop:
		for _, iv := range typedInput {
			mv, ok := iv.(map[string]interface{})
			if !ok {
				break loop
			}

			for k, v := range mv {
				if _, ok := f.children[k]; !ok {
					f.children[k] = newNode(k)
				}
				f.children[k].grow(v)
			}

		}

	}
}

func (f *node) repr() (string, error) {
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
		return b.String(), errors.New("invalid type definition")
	}

	b.Reset()
	printer.Fprint(&b, fset, astf)

	// remove go file header
	repr := b.String()
	repr = strings.TrimPrefix(repr, fileHeader)
	repr = strings.TrimSpace(repr)

	return repr, nil
}

func (f *node) reprToBuffer(b *bytes.Buffer) {
	if f.root {
		b.WriteString(fmt.Sprintf("type %s ", baseTypeName))
	} else {
		b.WriteString(fmt.Sprintf("%s ", attrName(f.name)))
	}
	b.WriteString(f.t.repr())

	isObject := false
	if f.t.id == nodeTypeObject || f.t.id == nodeTypeArrayObject {
		isObject = true

		b.WriteString(" {")
		defer b.WriteString("}")
	}

	if isObject && len(f.children) > 0 {
		// sort children by name
		type nodeWithName struct {
			name string
			node *node
		}
		var sortedChildren []nodeWithName
		for subName, child := range f.children {
			sortedChildren = append(sortedChildren, nodeWithName{
				name: subName,
				node: child,
			})
		}
		sort.Slice(sortedChildren, func(i, j int) bool {
			return sortedChildren[i].name < sortedChildren[j].name
		})

		// repr children
		for _, child := range sortedChildren {
			b.WriteString("\n")

			// repr child
			child.node.reprToBuffer(b)

			// add json tag
			b.WriteString(fmt.Sprintf(" `json:\"%s\"`", child.name))
		}
		b.WriteString("\n")
	}
}
