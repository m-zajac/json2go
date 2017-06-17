package json2go

import (
	"bytes"
	"fmt"
)

const baseTypeName = "Object"

type node struct {
	root     bool
	name     string
	t        nodeType
	required bool
	children map[string]*node
}

func newNode(name string) *node {
	return &node{
		name:     name,
		t:        newInitType(),
		required: true,
		children: make(map[string]*node),
	}
}

func (n *node) grow(input interface{}) {
	if input == nil {
		n.required = false
		return
	}

	keyFirstAppereance := (n.t.id == nodeTypeInit)
	growChild := func(k string, v interface{}, usedKeys map[string]struct{}) {
		if _, ok := n.children[k]; !ok {
			n.children[k] = newNode(k)
			if !keyFirstAppereance {
				n.children[k].required = false
			}
		}
		n.children[k].grow(v)
		usedKeys[k] = struct{}{}
	}

	n.t = n.t.grow(input)

	switch typedInput := input.(type) {
	case map[string]interface{}:
		usedKeys := make(map[string]struct{})
		for k, v := range typedInput {
			growChild(k, v, usedKeys)
		}

		for k, child := range n.children {
			if _, used := usedKeys[k]; !used {
				child.required = false
			}
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

			usedKeys := make(map[string]struct{})
			for k, v := range mv {
				growChild(k, v, usedKeys)
			}

			for k, child := range n.children {
				if _, used := usedKeys[k]; !used {
					child.required = false
				}
			}

			// after first iteration, no key is appearing as first
			// TODO: there should be a better way to implement this
			keyFirstAppereance = false
		}
	}
}

func (n *node) compare(n2 *node) bool {
	if n.name != n2.name {
		return false
	}
	if n.t.id != n2.t.id {
		return false
	}
	if n.required != n2.required {
		return false
	}
	if len(n.children) != len(n2.children) {
		return false
	}

	for k, child := range n.children {
		child2, ok := n2.children[k]
		if !ok {
			return false
		}
		if !child.compare(child2) {
			return false
		}
	}

	return true
}

func (n *node) repr(prefix string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s{\n", prefix))
	buf.WriteString(fmt.Sprintf("%s  name: %s\n", prefix, n.name))
	buf.WriteString(fmt.Sprintf("%s  type: %s\n", prefix, n.t.id))
	buf.WriteString(fmt.Sprintf("%s  required: %t\n", prefix, n.required))
	if len(n.children) > 0 {
		buf.WriteString(fmt.Sprintf("%s  children: {\n", prefix))
		for _, c := range n.children {
			buf.WriteString(fmt.Sprintf("%s    %s:\n%s\n", prefix, c.name, c.repr(prefix+"    ")))
		}
		buf.WriteString(fmt.Sprintf("%s  }\n", prefix))
	}
	buf.WriteString(fmt.Sprintf("%s}", prefix))

	return buf.String()
}
