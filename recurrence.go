package json2go

const (
	minNumOfAttributesForRecurrence = 3
)

func replaceRecurrentNodes(root *node) bool {
	return replaceRecurrent(root, nil)
}

func replaceRecurrent(n *node, parents []*node) bool {
	if n.t != nodeTypeObject {
		return false
	}
	if len(n.children) < minNumOfAttributesForRecurrence {
		return false
	}

	result := false
	parents = append(parents, n)
	for i, c := range n.children {
		if c.t != nodeTypeObject {
			continue
		}

		for _, p := range parents {
			if p.fits(c) {
				n.children[i] = &node{
					root:           false,
					nullable:       true, // recurrent type hase to be a pointer
					required:       c.required,
					key:            c.key,
					name:           c.name,
					t:              nodeTypeExtracted,
					externalTypeID: p.name,
					arrayLevel:     c.arrayLevel,
					arrayWithNulls: c.arrayWithNulls,
				}
				result = true
				break
			}
		}

		if replaceRecurrent(c, parents) {
			result = true
		}
	}

	return result
}
