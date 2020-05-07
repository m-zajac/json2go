package json2go

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	baseTypeName           = "Document"
	structIDlevelSeparator = "|"
)

type node struct {
	root           bool
	nullable       bool
	required       bool
	key            string
	name           string
	t              nodeType
	externalTypeID string
	children       []*node
	arrayLevel     int
	mapLevel       int
}

func newNode(key string) *node {
	return &node{
		key:      key,
		name:     attrName(key),
		t:        nodeTypeInit,
		nullable: false,
		required: true,
	}
}

func (n *node) grow(input interface{}) {
	if input == nil {
		n.nullable = true
		return
	}

	if n.t.id() == nodeTypeInterface.id() {
		return //nothing to do now
	}

	n.growChildrenFromData(input)

	switch typedInput := input.(type) {
	case []interface{}:
		if n.t != nodeTypeInit && n.arrayLevel == 0 {
			n.t = nodeTypeInterface
			n.children = nil
			break
		}

		localLevel, localType := arrayStructure(typedInput, n.t)
		if n.t == nodeTypeInit {
			n.t = localType
			n.arrayLevel = localLevel
		} else if n.arrayLevel != localLevel || n.t != localType {
			n.t = nodeTypeInterface
			n.arrayLevel = 0
		}
	default:
		n.t = growType(n.t, typedInput)
		n.arrayLevel = 0
	}
}

func (n *node) getOrCreateChild(key string) (*node, bool) {
	if child := n.getChild(key); child != nil {
		return child, false
	}

	childrenNames := make(map[string]bool)
	for _, c := range n.children {
		childrenNames[c.name] = true
	}

	child := newNode(key)

	for childrenNames[child.name] {
		child.name = nextName(child.name)
	}

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

func (n *node) growChildrenFromData(in interface{}) {
	if n.t == nodeTypeInterface {
		return
	}

	if ar, ok := in.([]interface{}); ok {
		for i := range ar {
			n.growChildrenFromData(ar[i])
		}
		return
	}

	obj, ok := in.(map[string]interface{})
	if !ok {
		n.children = nil
		return
	}

	alreadyHasChildren := (n.children != nil)
	usedKeys := make(map[string]bool)
	for k, v := range obj {
		child, created := n.getOrCreateChild(k)
		if created && alreadyHasChildren {
			child.required = false
		}
		child.grow(v)
		usedKeys[k] = true
	}

	for _, child := range n.children {
		if !usedKeys[child.key] {
			child.required = false
		}
	}
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
	if n.t.id() != n2.t.id() {
		return false
	}
	if n.nullable != n2.nullable {
		return false
	}
	if n.required != n2.required {
		return false
	}
	if n.externalTypeID != n2.externalTypeID {
		return false
	}
	if n.arrayLevel != n2.arrayLevel {
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
	if n.arrayLevel > 0 {
		buf.WriteString(fmt.Sprintf("%s  type: [%d]%s\n", prefix, n.arrayLevel, n.t.id()))
	} else {
		buf.WriteString(fmt.Sprintf("%s  type: %s\n", prefix, n.t.id()))
	}
	buf.WriteString(fmt.Sprintf("%s  nullable: %t\n", prefix, n.nullable))
	buf.WriteString(fmt.Sprintf("%s  required: %t\n", prefix, n.required))
	if n.externalTypeID != "" {
		buf.WriteString(fmt.Sprintf("%s  extType: %s\n", prefix, n.externalTypeID))
	}
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

// arrayStructure returns array depth and elements type. If array is nested and has no consistent structure, level -1 is returned.
func arrayStructure(in []interface{}, inType nodeType) (int, nodeType) {
	if inType == nil {
		inType = nodeTypeInit
	}
	if len(in) == 0 {
		return 1, inType
	}

	depth := 0
	for _, el := range in {
		switch typedEl := el.(type) {
		case []interface{}:
			localDepth, localType := arrayStructure(typedEl, inType)
			localDepth++

			if inType == nodeTypeInit {
				inType = localType
			} else if localType != inType {
				if localType.expands(inType) {
					inType = localType
				} else {
					inType = nodeTypeInterface
				}
			}

			switch depth {
			case 0:
				depth = localDepth
			case localDepth:
			default:
				return -1, nodeTypeInterface
			}
		default:
			localType := inType.fit(typedEl)

			if inType == nodeTypeInit {
				inType = localType
			} else if localType != inType {
				if localType.expands(inType) {
					inType = localType
				} else {
					inType = nodeTypeInterface
				}
			}

			depth = 1
		}
	}

	return depth, inType
}
