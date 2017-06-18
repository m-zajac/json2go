package json2go

import (
	"bytes"
	"fmt"
	"sort"
)

const baseTypeName = "Object"

type node struct {
	root     bool
	key      string
	t        nodeType
	required bool
	children []*node
}

func newNode(name string) *node {
	return &node{
		key:      name,
		t:        newInitType(),
		required: true,
	}
}

func (n *node) grow(input interface{}) {
	if input == nil {
		n.required = false
		return
	}

	keyFirstAppereance := (n.t.id == nodeTypeInit)
	growChild := func(k string, v interface{}, usedKeys map[string]struct{}) {
		child, created := n.getOrCreateChild(k)
		if created && !keyFirstAppereance {
			child.required = false
		}
		child.grow(v)
		usedKeys[k] = struct{}{}
	}

	n.t = n.t.grow(input)

	switch typedInput := input.(type) {
	case map[string]interface{}:
		usedKeys := make(map[string]struct{})
		for k, v := range typedInput {
			growChild(k, v, usedKeys)
		}

		for _, child := range n.children {
			if _, used := usedKeys[child.key]; !used {
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

			for _, child := range n.children {
				if _, used := usedKeys[child.key]; !used {
					child.required = false
				}
			}

			// after first iteration, no key is appearing as first
			// TODO: there should be a better way to implement this
			keyFirstAppereance = false
		}
	}
}

func (n *node) getOrCreateChild(key string) (*node, bool) {
	if child := n.getChild(key); child != nil {
		return child, false
	}

	child := newNode(key)
	n.children = append(n.children, child)
	return child, true
}

func (n *node) getChild(key string) *node {
	for _, child := range n.children {
		if child.key == key {
			return child
		}
	}

	return nil
}

func (n *node) sort() {
	sort.Slice(n.children, func(i int, j int) bool {
		return n.children[i].key < n.children[j].key
	})

	for _, child := range n.children {
		child.sort()
	}
}

func (n *node) compare(n2 *node) bool {
	if n.key != n2.key {
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

	for i, child := range n.children {
		child2 := n2.children[i]
		if !child.compare(child2) {
			return false
		}
	}

	return true
}

func (n *node) repr(prefix string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s{\n", prefix))
	buf.WriteString(fmt.Sprintf("%s  key: %s\n", prefix, n.key))
	buf.WriteString(fmt.Sprintf("%s  type: %s\n", prefix, n.t.id))
	buf.WriteString(fmt.Sprintf("%s  required: %t\n", prefix, n.required))
	if len(n.children) > 0 {
		buf.WriteString(fmt.Sprintf("%s  children: {\n", prefix))
		for _, c := range n.children {
			buf.WriteString(fmt.Sprintf("%s    %s:\n%s\n", prefix, c.key, c.repr(prefix+"    ")))
		}
		buf.WriteString(fmt.Sprintf("%s  }\n", prefix))
	}
	buf.WriteString(fmt.Sprintf("%s}", prefix))

	return buf.String()
}
