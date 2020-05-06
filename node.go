package json2go

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

const (
	baseTypeName           = "Document"
	structIDlevelSeparator = "|"
)

type node struct {
	root           bool
	nullable       bool
	required       bool
	arrayLevel     int
	key            string
	name           string
	t              nodeType
	externalTypeID string
	children       []*node
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

// structureID returns identifier unique for this nodes structure
// if `asRoot` is true, this node id does not depend on "nullable" or "required" property
func (n *node) structureID(asRoot bool) string {
	var id string
	if asRoot {
		id = n.t.id()
	} else {
		id = fmt.Sprintf("%s.%s", n.key, n.t.id())
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
	typeID      string
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
			typeID:      n.t.id(),
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

// extractCommonSubtree extracts at most one common subtree to new root node
func extractCommonSubtree(root *node, rootNames map[string]bool) *node {
	infos := make(map[string]nodeStructureInfo)

	root.treeInfo(infos)

	// Find all objects with at least 2 children nodes.
	infosForExtraction := make([]nodeStructureInfo, 0, len(infos))
	for _, info := range infos {
		if len(info.nodes) > 1 && info.typeID == nodeTypeObject.id() {
			infosForExtraction = append(infosForExtraction, info)
		}
	}
	if len(infosForExtraction) == 0 {
		return nil
	}

	// Sort by tree depth, asceding. We want to start extracting from simpliest subtrees.
	sort.Slice(infosForExtraction, func(i int, j int) bool {
		l1 := strings.Count(infosForExtraction[i].structureID, structIDlevelSeparator)
		l2 := strings.Count(infosForExtraction[j].structureID, structIDlevelSeparator)

		if l1 == l2 { // if struct depth is equal, compare by first node key
			return infosForExtraction[i].nodes[0].key < infosForExtraction[j].nodes[0].key
		}

		// Compare by struct depth.
		return l1 < l2
	})

	for _, info := range infosForExtraction {
		extractedKey, extractedName := makeNameFromNodes(info.nodes)
		if extractedName == "" {
			continue
		}

		for rootNames[extractedName] {
			extractedName = nextName(extractedName)
			extractedKey = nextName(extractedKey)
		}
		rootNames[extractedName] = true

		extractedNode := mergeNodes(info.nodes)
		extractedNode.name = extractedName
		extractedNode.key = extractedKey
		extractedNode.root = true

		root.modify(info.structureID, func(modNode *node) {
			modNode.t = nodeTypeExtracted
			modNode.externalTypeID = extractedName
			modNode.children = nil
		})

		return extractedNode // exit after first successful extract
	}

	return nil
}

func extractCommonSubtrees(root *node) []*node {
	rootNames := map[string]bool{
		root.name: true,
	}

	extractedSize := 0
	nodes := []*node{root}
	for len(nodes) != extractedSize {
		extractedSize = len(nodes)
		result := nodes
		for _, n := range nodes {
			extNode := extractCommonSubtree(n, rootNames)
			if extNode != nil {
				result = append(result, extNode)
			}
		}
		nodes = result
	}

	return nodes
}

// makeNameFromNodes is helper function trying to find the best name (and key) from list of nodes.
func makeNameFromNodes(nodes []*node) (key, name string) {
	if len(nodes) == 0 {
		return "", ""
	}

	// Try to create name from all node names.
	var names []string
	var keys []string
	for _, in := range nodes {
		names = append(names, in.name)
		keys = append(keys, in.key)
	}
	name = extractCommonName(names...)
	key = extractCommonName(keys...)

	// If unsuccessful, try to create name from first node's attribute names.
	if name == "" {
		var keys []string
		for _, child := range nodes[0].children {
			keys = append(keys, child.key)
		}
		key = nameFromNames(keys...)
		name = attrName(key)
	}

	return key, name
}

// mergeNodes merges multiple nodes into one.
// If any of nodes is not required, merged node is also not required.
// If any of nodes is nullable, merged node is also nullable.
// Children of merged nodes are also merged by the same rules.
func mergeNodes(nodes []*node) *node {
	if len(nodes) == 0 { // This should never happen
		panic("mergeNodes called with empty node list!")
	}
	if len(nodes) == 1 {
		return nodes[0]
	}

	// Set main attributes of merged node.
	merged := *nodes[0]
	for _, n := range nodes {
		if !n.required {
			merged.required = false
		}
		if n.nullable {
			merged.nullable = true
		}
	}

	// Set attributes of merged node's children recurently.
	if len(merged.children) > 0 {
		for i, cn := range merged.children {
			cnodes := make([]*node, 0, len(nodes))
			for _, n := range nodes {
				v := n.getChild(cn.key)
				if v == nil {
					continue
				}
				cnodes = append(cnodes, v)
			}
			if len(cnodes) > 1 {
				merged.children[i] = mergeNodes(cnodes)
			}
		}
	}

	return &merged
}
