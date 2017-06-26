package json2go

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

const (
	baseTypeName           = "Object"
	structIDlevelSeparator = "|"
)

type node struct {
	root           bool
	key            string
	t              nodeType
	externalTypeID string
	required       bool
	children       []*node
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
	if n.externalTypeID != n2.externalTypeID {
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

// structureID returns identifier unique for this nodes structure
// if `asRoot` is true, this node id does not depend on "required" property
func (n *node) structureID(asRoot bool) string {
	var id string
	if asRoot {
		id = string(n.t.id)
	} else {
		id = fmt.Sprintf("%s.%s.%t", n.key, n.t.id, n.required)
	}

	var parts []string
	for _, child := range n.children {
		parts = append(parts, child.structureID(false))
	}

	result := id
	if len(parts) > 0 {
		result += structIDlevelSeparator + strings.Join(parts, ",")
	}
	return result
}

type nodeStructureInfo struct {
	structureID string
	typeID      nodeTypeID
	nodes       []*node
}

func (n *node) treeInfo(infos map[string]nodeStructureInfo) {
	var info nodeStructureInfo

	id := n.structureID(true)
	if ninfo, ok := infos[id]; ok {
		info = ninfo
		info.nodes = append(info.nodes, n)
	} else {
		info = nodeStructureInfo{
			structureID: id,
			typeID:      n.t.id,
			nodes:       []*node{n},
		}
	}
	infos[id] = info

	for _, child := range n.children {
		child.treeInfo(infos)
	}
}

// modify executes function f on all nodes in subtree with given structure id
func (n *node) modify(structureID string, f func(*node)) {
	for i, child := range n.children {
		if child.structureID(true) == structureID {
			f(n.children[i])
		}

		child.modify(structureID, f)
	}
}

func (n *node) repr(prefix string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s{\n", prefix))
	buf.WriteString(fmt.Sprintf("%s  key: %s\n", prefix, n.key))
	buf.WriteString(fmt.Sprintf("%s  type: %s\n", prefix, n.t.id))
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

// extractCommonSubtree extracts at most one common subtree to new root node
func extractCommonSubtree(root *node, rootKeys map[string]struct{}) *node {
	infos := make(map[string]nodeStructureInfo)

	root.treeInfo(infos)

	infosForExtraction := make([]nodeStructureInfo, 0, len(infos))
	for _, info := range infos {
		if len(info.nodes) > 1 && info.typeID == nodeTypeObject {
			infosForExtraction = append(infosForExtraction, info)
		}
	}
	sort.Slice(infosForExtraction, func(i int, j int) bool {
		l1 := strings.Count(infosForExtraction[i].structureID, structIDlevelSeparator)
		l2 := strings.Count(infosForExtraction[j].structureID, structIDlevelSeparator)

		if l1 == l2 { // if struct depth is equal, compare by first node key
			return infosForExtraction[i].nodes[0].key < infosForExtraction[j].nodes[0].key
		}

		// compare by struct depth
		return l1 < l2
	})

	for _, info := range infosForExtraction {
		extractedNode := *info.nodes[0]

		var names []string
		for _, in := range info.nodes {
			names = append(names, in.key)
		}
		extractedKey := extractCommonName(names...)
		if extractedKey == "" {
			var keys []string
			for _, child := range extractedNode.children {
				keys = append(keys, child.key)
			}
			extractedKey = keynameFromKeys(keys...)
		}
		if extractedKey == "" {
			continue
		}

		_, exists := rootKeys[extractedKey]
		for exists {
			extractedKey = nextName(extractedKey)
			_, exists = rootKeys[extractedKey]
		}
		rootKeys[extractedKey] = struct{}{}
		extractedNode.key = extractedKey
		extractedNode.root = true

		root.modify(info.structureID, func(modNode *node) {
			modNode.t = newExternalObjectType()
			modNode.externalTypeID = attrName(extractedKey)
			modNode.children = nil
		})

		return &extractedNode // exit after first successful extract
	}

	return nil
}

func extractCommonSubtrees(root *node) []*node {
	rootKeys := map[string]struct{}{
		root.key: struct{}{},
	}

	extractedSize := 0
	nodes := []*node{root}
	for len(nodes) != extractedSize {
		extractedSize = len(nodes)
		result := nodes
		for _, n := range nodes {
			extNode := extractCommonSubtree(n, rootKeys)
			if extNode != nil {
				result = append(result, extNode)
			}
		}
		nodes = result
	}

	return nodes
}
