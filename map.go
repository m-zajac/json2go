package json2go

import "strings"

func convertViableObjectsToMaps(root *node, minAttributes uint) {
	for _, c := range root.children {
		if c.t.id() != nodeTypeObject.id() {
			continue
		}

		convertViableObjectsToMaps(c, minAttributes)

		if tryConvertToMap(c, minAttributes) {
			continue
		}
	}

	tryConvertToMap(root, minAttributes)
}

func tryConvertToMap(n *node, minAttributes uint) bool {
	if len(n.children) < int(minAttributes) {
		return false
	}
	if len(n.children) < 1 {
		return false
	}

	// Children has to have same type and structure.
	t := n.children[0].t
	sid := mergeNumsStructureID(n.children[0], true)
	for _, c := range n.children {
		if !t.expands(c.t) && !c.t.expands(t) {
			return false
		}
		if mergeNumsStructureID(c, true) != sid {
			return false
		}
	}

	// Convert this node to map.
	newNode := mergeNodes(n.children)
	newNode.mapLevel = n.children[0].mapLevel + 1
	newNode.key = n.key
	newNode.name = n.name
	newNode.arrayLevel = n.arrayLevel
	newNode.root = n.root
	*n = *newNode

	return true
}

func mergeNumsStructureID(n *node, asRoot bool) string {
	sid := structureID(n, asRoot)
	sid = strings.Replace(sid, nodeTypeInt.id(), "number", -1)
	sid = strings.Replace(sid, nodeTypeFloat.id(), "number", -1)

	return sid
}
