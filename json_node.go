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

func (n *node) grow(input interface{}) {
	n.t = n.t.grow(input)

	switch typedInput := input.(type) {
	case map[string]interface{}:
		for k, v := range typedInput {
			if _, ok := n.children[k]; !ok {
				n.children[k] = newNode(k)
			}
			n.children[k].grow(v)
		}
	case []interface{}:
		if n.t.id != nodeTypeArrayObject {
			break
		}

	loop:
		for _, iv := range typedInput {
			mv, ok := iv.(map[string]interface{})
			if !ok {
				break loop
			}

			for k, v := range mv {
				if _, ok := n.children[k]; !ok {
					n.children[k] = newNode(k)
				}
				n.children[k].grow(v)
			}

		}

	}
}

func (n *node) repr() (string, error) {
	var b bytes.Buffer

	// write valid go file header
	fileHeader := "package main\n\n"
	b.WriteString(fileHeader)

	// write struct representation
	if !n.reprToBuffer(&b) {
		return "", errors.New("cannot produce valid type definition")
	}

	// parse and print formatted
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "", b.String(), 0)
	if err != nil {
		return b.String(), errors.New("cannot produce valid type definition")
	}

	b.Reset()
	printer.Fprint(&b, fset, astf)

	// remove go file header
	repr := b.String()
	repr = strings.TrimPrefix(repr, fileHeader)
	repr = strings.TrimSpace(repr)

	return repr, nil
}

func (n *node) reprToBuffer(b *bytes.Buffer) bool {
	if n.root {
		b.WriteString(fmt.Sprintf("type %s ", baseTypeName))
	} else {
		name := attrName(n.name)
		if name == "" {
			return false
		}
		b.WriteString(fmt.Sprintf("%s ", name))
	}
	b.WriteString(n.t.repr())

	isObject := false
	if n.t.id == nodeTypeObject || n.t.id == nodeTypeArrayObject {
		isObject = true

		b.WriteString(" {")
		defer b.WriteString("}")
	}

	if isObject && len(n.children) > 0 {
		// sort children by name
		type nodeWithName struct {
			name string
			node *node
		}
		var sortedChildren []nodeWithName
		for subName, child := range n.children {
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
			if !child.node.reprToBuffer(b) {
				continue
			}

			// add json tag
			b.WriteString(fmt.Sprintf(" `json:\"%s\"`", child.name))
		}
		b.WriteString("\n")
	}

	return true
}
