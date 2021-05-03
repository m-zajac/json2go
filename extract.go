package json2go

import (
	"fmt"
	"sort"
	"strings"
)

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

// extractCommonSubtree extracts at most one common subtree to new root node
func extractCommonSubtree(root *node, rootNames map[string]bool) *node {
	// Find all structures in object tree.
	structDataM := make(map[string]structNodes)
	objectTreeInfo(root, structDataM)

	// Filter out structures that shouldn't be extracted.
	var keysToDel []string
	for k, info := range structDataM {
		if len(info.nodes) < 2 {
			// This structure occurs only once in document, so nothing to extract.
			keysToDel = append(keysToDel, k)
		} else if len(info.nodes) > 0 && info.nodes[0].t.id() == nodeTypeMap.id() {
			// Don't extract maps!
			keysToDel = append(keysToDel, k)
		}
	}
	for _, k := range keysToDel {
		delete(structDataM, k)
	}
	if len(structDataM) == 0 {
		return nil
	}

	// Create list, sorted by tree depth, asceding. We want to start extracting from simpliest subtrees.
	structData := make([]structNodes, 0, len(structDataM))
	for _, v := range structDataM {
		structData = append(structData, v)
	}
	sort.Slice(structData, func(i int, j int) bool {
		l1 := strings.Count(structData[i].structureID, structIDlevelSeparator)
		l2 := strings.Count(structData[j].structureID, structIDlevelSeparator)

		if l1 == l2 { // if struct depth is equal, compare by first node key
			return structData[i].nodes[0].key < structData[j].nodes[0].key
		}
		return l1 < l2
	})

	for _, info := range structData {
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
		extractedNode.arrayLevel = 0

		modifyTree(root, info.structureID, func(modNode *node) {
			modNode.t = nodeTypeExtracted
			modNode.externalTypeID = extractedName
			modNode.children = nil
		})

		return extractedNode // exit after first successful extract
	}

	return nil
}

type structNodes struct {
	structureID string
	nodes       []*node
}

func objectTreeInfo(n *node, infos map[string]structNodes) {
	switch n.t.id() {
	case nodeTypeObject.id():
	case nodeTypeMap.id():
	default:
		return
	}

	var info structNodes

	id := structureID(n, false)
	if ninfo, ok := infos[id]; ok {
		info = ninfo
		info.nodes = append(info.nodes, n)
	} else {
		info = structNodes{
			structureID: id,
			nodes:       []*node{n},
		}
	}
	infos[id] = info

	for _, child := range n.children {
		objectTreeInfo(child, infos)
	}
}

// structureID returns identifier unique for this nodes structure
// if `withKey` is true, node's key name is added to id.
func structureID(n *node, withKey bool) string {
	id := n.t.id()
	if withKey {
		id = fmt.Sprintf("%s.%s", n.key, id)
	}

	var parts []string
	for _, child := range n.children {
		parts = append(parts, structureID(child, true))
	}

	result := id
	if len(parts) > 0 {
		result += structIDlevelSeparator + strings.Join(parts, ",")
	}
	return result
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
		if n.t.expands(merged.t) {
			merged.t = n.t
		}

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

// modifyTree executes function f on all nodes in subtree with given structure id
func modifyTree(root *node, structID string, f func(*node)) {
	for i, child := range root.children {
		if structureID(child, false) == structID {
			f(root.children[i])
		}

		modifyTree(child, structID, f)
	}
}
