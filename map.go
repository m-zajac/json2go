package json2go

import "strings"

func convertViableObjectsToMaps(root *node, minAttributes uint) {
	// for _, c := range root.children {
	// 	if c.t.id() != nodeTypeObject.id() {
	// 		continue
	// 	}

	// 	convertViableObjectsToMaps(c, minAttributes)

	// 	if tryConvertToMap(c, minAttributes) {
	// 		continue
	// 	}
	// }

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
	sid := mergeNumsStructureID(n.children[0], false)
	for _, c := range n.children {
		if !t.expands(c.t) && !c.t.expands(t) {
			return false
		}
		if mergeNumsStructureID(c, false) != sid {
			return false
		}
	}

	// Convert this node to map.
	n.t = nodeTypeMap

	// Add child as map value type node
	newNode := mergeNodes(n.children)
	newNode.key = ""
	newNode.name = ""
	newNode.required = true
	n.children = []*node{newNode}

	return true
}

func mergeNumsStructureID(n *node, withKey bool) string {
	sid := structureID(n, withKey)
	sid = strings.Replace(sid, nodeTypeInt.id(), "number", -1)
	sid = strings.Replace(sid, nodeTypeFloat.id(), "number", -1)

	return sid
}
