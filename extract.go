package json2go

import (
	"fmt"
	"sort"
	"strings"
)

func extractCommonSubtrees(root *node, opts options) []*node {
	rootNames := map[string]bool{
		root.name: true,
	}

	nodes := []*node{root}
	for {
		changed := false
		extNode := extractOne(root, nodes, rootNames, opts, &changed)
		if extNode != nil {
			nodes = append(nodes, extNode)
			changed = true
		}
		if !changed {
			break
		}
	}

	// Final pass: promote any anonymous structs with embeddings to named types.
	for {
		promoted := false
		var candidates []*node
		for _, r := range nodes {
			collectAnonymousEmbeddedNodes(r, &candidates)
		}

		if len(candidates) > 0 {
			// Extract the first candidate
			c := candidates[0]
			_, name := makeNameFromNodes([]*node{c}, nil)
			for rootNames[name] {
				name = nextName(name)
			}
			rootNames[name] = true

			extNode := c.clone()
			extNode.name = name
			extNode.root = true
			extNode.arrayLevel = 0

			c.t = nodeTypeExtracted
			c.externalTypeID = name
			c.children = nil

			nodes = append(nodes, extNode)
			promoted = true
		}

		if !promoted {
			break
		}
	}

	return nodes
}

func collectAnonymousEmbeddedNodes(n *node, candidates *[]*node) {
	if n.t.id() == nodeTypeObject.id() && !n.root && n.externalTypeID == "" {
		hasEmbedding := false
		for _, child := range n.children {
			if child.embedded {
				hasEmbedding = true
				break
			}
		}
		if hasEmbedding {
			*candidates = append(*candidates, n)
			return // Don't look deeper if we found a candidate
		}
	}

	for _, child := range n.children {
		collectAnonymousEmbeddedNodes(child, candidates)
	}
}

func extractOne(root *node, allRoots []*node, rootNames map[string]bool, opts options, changed *bool) *node {
	// 1. Collect all object/map nodes from all current roots.
	var candidates []*node
	parentKeys := make(map[*node][]string)
	for _, r := range allRoots {
		collectObjectNodes(r, &candidates, parentKeys)
	}

	// 2. Try to find identity/similarity-based merges.
	if ext := tryExtractSimilarity(root, candidates, allRoots, rootNames, opts, changed, parentKeys); ext != nil || *changed {
		return ext
	}

	// 3. Try to find common subsets for embedding.
	if ext := tryExtractSubset(root, candidates, allRoots, rootNames, opts, changed, parentKeys); ext != nil || *changed {
		return ext
	}

	return nil
}

func collectObjectNodes(n *node, candidates *[]*node, parentKeys map[*node][]string) {
	if n.t.id() == nodeTypeObject.id() {
		*candidates = append(*candidates, n)
	}

	for _, child := range n.children {
		if child.t.id() == nodeTypeObject.id() {
			parentKeys[child] = append(parentKeys[child], child.key)
		}
		collectObjectNodes(child, candidates, parentKeys)
	}
}

func tryExtractSimilarity(root *node, candidates []*node, allRoots []*node, rootNames map[string]bool, opts options, changed *bool, parentKeys map[*node][]string) *node {
	// Prioritize recursion: check if any candidate is similar to an existing root.
	for _, c := range candidates {
		if c.root {
			continue // Already a root
		}

		for _, r := range allRoots {
			// If it's the root node, we want to be very sure it's recursive.
			threshold := opts.extractSimilarityThreshold
			if r.root && r.key == root.key {
				threshold = 0.9 // Higher threshold for root node
			}

			sim := c.similarity(r)
			if sim >= threshold {
				// Found a match!
				extractedName := r.externalTypeID
				if extractedName == "" {
					extractedName = r.name
				}

				c.t = nodeTypeExtracted
				c.externalTypeID = extractedName
				c.children = nil
				*changed = true
				return nil
			}
		}
	}

	// Then check similarity between candidates.
	for i, c1 := range candidates {
		if c1.root {
			continue
		}
		var matches []*node
		for j := i + 1; j < len(candidates); j++ {
			c2 := candidates[j]
			if c2.root {
				continue
			}
			if c1.similarity(c2) >= opts.extractSimilarityThreshold {
				matches = append(matches, c2)
			}
		}

		if len(matches) > 0 {
			matches = append(matches, c1)
			*changed = true
			return extractNodes(root, matches, rootNames, parentKeys)
		}
	}

	return nil
}

func extractNodes(root *node, nodes []*node, rootNames map[string]bool, parentKeys map[*node][]string) *node {
	var allParentKeys []string
	for _, n := range nodes {
		allParentKeys = append(allParentKeys, parentKeys[n]...)
	}
	extractedKey, extractedName := makeNameFromNodes(nodes, allParentKeys)
	if extractedName == "" {
		extractedName = "Extracted"
		extractedKey = "extracted"
	}
	for rootNames[extractedName] {
		extractedName = nextName(extractedName)
		extractedKey = nextName(extractedKey)
	}
	rootNames[extractedName] = true

	extractedNode := mergeNodes(nodes)
	extractedNode.name = extractedName
	extractedNode.key = extractedKey
	extractedNode.root = true
	extractedNode.arrayLevel = 0

	for _, n := range nodes {
		n.t = nodeTypeExtracted
		n.externalTypeID = extractedName
		n.children = nil
	}

	return extractedNode
}

func tryExtractSubset(root *node, candidates []*node, allRoots []*node, rootNames map[string]bool, opts options, changed *bool, parentKeys map[*node][]string) *node {
	if opts.extractMinSubsetSize < 1 {
		return nil
	}

	// Count occurrences of subsets of keys+types.
	// For simplicity, we'll look for pairs of nodes and their common fields.
	type subsetMatch struct {
		nodes []*node
		keys  []string
	}
	var best subsetMatch

	for i, c1 := range candidates {
		for j := i + 1; j < len(candidates); j++ {
			c2 := candidates[j]
			common := getCommonFields(c1, c2)
			if len(common) >= opts.extractMinSubsetSize {
				// Check if this subset is exactly one of the current roots.
				var existingRoot *node
				for _, r := range allRoots {
					if r.t.id() == nodeTypeObject.id() && len(r.children) == len(common) {
						rFields := getCommonFields(r, r) // Get all fields of r
						if hasAllFields(r, common) && len(rFields) == len(common) {
							existingRoot = r
							break
						}
					}
				}

				if existingRoot != nil {
					// Instead of extracting a new subset, just use the existing root in candidates that have these fields.
					matchFound := false
					for _, c3 := range candidates {
						if c3.root || !hasAllFields(c3, common) {
							continue
						}
						// Check if c3 already embeds this root.
						alreadyEmbeds := false
						for _, child := range c3.children {
							if child.embedded && child.externalTypeID == existingRoot.name {
								alreadyEmbeds = true
								break
							}
						}
						if alreadyEmbeds {
							continue
						}

						if collidesWithFields(existingRoot.name, c3, common) {
							continue
						}

						// Replace fields with embedding.
						newChildren := make([]*node, 0, len(c3.children)-len(common)+1)
						newChildren = append(newChildren, &node{
							t:              nodeTypeExtracted,
							externalTypeID: existingRoot.name,
							embedded:       true,
							required:       true,
						})
						for _, child := range c3.children {
							found := false
							for _, k := range common {
								if child.key == k {
									found = true
									break
								}
							}
							if !found {
								newChildren = append(newChildren, child)
							}
						}
						c3.children = newChildren
						matchFound = true
					}
					if matchFound {
						*changed = true
						return nil
					}
					continue
				}

				// Found a potential subset.
				// Count how many candidates have all these fields.
				var matches []*node
				for _, c3 := range candidates {
					if hasAllFields(c3, common) {
						matches = append(matches, c3)
					}
				}

				if len(matches) >= opts.extractMinSubsetOccurrences {
					if len(common) > len(best.keys) || (len(common) == len(best.keys) && len(matches) > len(best.nodes)) {
						best = subsetMatch{nodes: matches, keys: common}
					}
				}
			}
		}
	}

	if len(best.nodes) > 0 {
		// Micro-embedding prevention:
		// If this subset would embed exactly ONE other extracted type AND add only a few fields, skip it.
		embeddedCount := 0
		for _, k := range best.keys {
			if strings.HasPrefix(k, "|") {
				embeddedCount++
			}
		}
		newFieldsCount := len(best.keys) - embeddedCount
		if embeddedCount > 0 && newFieldsCount < opts.extractMinAddedFields {
			return nil
		}

		// Extract the subset.
		var actualKeys []string
		for _, k := range best.keys {
			if !strings.HasPrefix(k, "|") {
				actualKeys = append(actualKeys, k)
			}
		}

		var allParentKeys []string
		for _, n := range best.nodes {
			allParentKeys = append(allParentKeys, parentKeys[n]...)
		}

		extractedKey, extractedName := makeNameFromNodes(nil, allParentKeys)
		if extractedName == "" || len(extractedName) <= 2 {
			extractedKey = nameFromNamesCapped(actualKeys...)
			extractedName = attrName(extractedKey)
		}

		if extractedName == "" {
			extractedName = "Extracted"
			extractedKey = "extracted"
		}
		for {
			if rootNames[extractedName] {
				extractedName = nextName(extractedName)
				extractedKey = nextName(extractedKey)
				continue
			}

			collides := false
			for _, n := range best.nodes {
				if len(n.children) == len(best.keys) {
					continue
				}
				if collidesWithFields(extractedName, n, best.keys) {
					collides = true
					break
				}
			}

			if collides {
				extractedName = nextName(extractedName)
				extractedKey = nextName(extractedKey)
				continue
			}

			break
		}
		rootNames[extractedName] = true

		// Create the base node.
		baseNode := newNode(extractedKey)
		baseNode.name = extractedName
		baseNode.t = nodeTypeObject
		baseNode.root = true

		for _, k := range best.keys {
			if strings.HasPrefix(k, "|") {
				extID := strings.TrimPrefix(k, "|")
				baseNode.children = append(baseNode.children, &node{
					t:              nodeTypeExtracted,
					externalTypeID: extID,
					embedded:       true,
					required:       true,
				})
				continue
			}

			var fieldInstances []*node
			for _, n := range best.nodes {
				if child := n.getChild(k); child != nil {
					fieldInstances = append(fieldInstances, child)
				}
			}
			if len(fieldInstances) > 0 {
				mergedField := mergeNodes(fieldInstances)
				mergedField.root = false
				baseNode.children = append(baseNode.children, mergedField)
			}
		}

		// Modify matching nodes to embed this base node.
		for _, n := range best.nodes {
			// If all fields were in the subset, replace node with a reference to the base node.
			if len(n.children) == len(best.keys) {
				n.t = nodeTypeExtracted
				n.externalTypeID = extractedName
				n.children = nil
				continue
			}

			// Remove fields from node and add embedded field.
			newChildren := make([]*node, 0, len(n.children)-len(best.keys)+1)
			embedded := &node{
				t:              nodeTypeExtracted,
				externalTypeID: extractedName,
				embedded:       true,
				required:       true, // Base structs are usually required if fields were there
			}
			newChildren = append(newChildren, embedded)
			for _, child := range n.children {
				found := false
				for _, k := range best.keys {
					if strings.HasPrefix(k, "|") {
						if child.embedded && child.externalTypeID == strings.TrimPrefix(k, "|") {
							found = true
							break
						}
					} else if child.key == k {
						found = true
						break
					}
				}
				if !found {
					newChildren = append(newChildren, child)
				}
			}
			n.children = newChildren
		}

		*changed = true
		return baseNode
	}

	return nil
}

func getCommonFields(n1, n2 *node) []string {
	set1 := make(map[string]string)
	for _, c := range n1.children {
		key := c.key
		if c.embedded {
			key = "|" + c.externalTypeID // Use externalTypeID for comparison of embedded fields
		}
		if key != "" {
			set1[key] = c.t.id()
		}
	}

	var common []string
	for _, c := range n2.children {
		key := c.key
		if c.embedded {
			key = "|" + c.externalTypeID
		}
		if key != "" && set1[key] == c.t.id() {
			common = append(common, key)
		}
	}
	sort.Strings(common)
	return common
}

func hasAllFields(n *node, keys []string) bool {
	for _, k := range keys {
		if strings.HasPrefix(k, "|") {
			extID := strings.TrimPrefix(k, "|")
			found := false
			for _, child := range n.children {
				if child.embedded && child.externalTypeID == extID {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		} else if n.getChild(k) == nil {
			return false
		}
	}
	return true
}

func makeNameFromNodes(nodes []*node, parentKeys []string) (key, name string) {
	// Try to find a common name from parent keys first.
	if len(parentKeys) > 0 {
		key = extractCommonName(parentKeys...)
		name = attrName(key)
		if name != "" && len(name) > 2 {
			return key, name
		}
	}

	if len(nodes) == 0 {
		return "", ""
	}

	// Try to create name from all node keys.
	var keys []string
	for _, in := range nodes {
		if in.key != "" {
			keys = append(keys, in.key)
		}
	}
	key = extractCommonName(keys...)
	name = attrName(key)

	// If successful and name is descriptive enough, return it.
	if name != "" && len(name) > 2 {
		return key, name
	}

	// If unsuccessful, try to create name from all node names.
	var names []string
	for _, in := range nodes {
		if in.name != "" {
			names = append(names, in.name)
		}
	}
	name = extractCommonName(names...)
	key = extractCommonName(keys...)

	// If still unsuccessful or too short, try to create name from first node's attribute names.
	if name == "" || len(name) <= 2 {
		var keys []string
		for _, child := range nodes[0].children {
			if child.key != "" {
				keys = append(keys, child.key)
			}
		}
		key = nameFromNamesCapped(keys...)
		name = attrName(key)
	}

	if name == "" {
		key = "extracted"
		name = "Extracted"
	}

	return key, name
}

func mergeNodes(nodes []*node) *node {
	if len(nodes) == 0 { // This should never happen
		panic("mergeNodes called with empty node list!")
	}
	if len(nodes) == 1 {
		return nodes[0].clone()
	}

	// Set main attributes of merged node.
	merged := nodes[0].clone()
	for _, n := range nodes[1:] {
		if n.t.expands(merged.t) {
			merged.t = n.t
		} else if !merged.t.expands(n.t) {
			merged.t = nodeTypeInterface
		}

		if !n.required {
			merged.required = false
		}
		if n.nullable {
			merged.nullable = true
		}
		if n.arrayWithNulls {
			merged.arrayWithNulls = true
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
			} else if len(cnodes) == 1 {
				merged.children[i] = cnodes[0].clone()
			}
		}
	}

	return merged
}

func isFieldExcluded(child *node, excludingKeys []string) bool {
	for _, k := range excludingKeys {
		if strings.HasPrefix(k, "|") {
			if child.embedded && child.externalTypeID == strings.TrimPrefix(k, "|") {
				return true
			}
		} else if child.key == k {
			return true
		}
	}
	return false
}

func collidesWithFields(name string, n *node, excludingKeys []string) bool {
	for _, child := range n.children {
		if isFieldExcluded(child, excludingKeys) {
			continue
		}
		if child.embedded {
			if child.externalTypeID == name {
				return true
			}
		} else if attrName(child.key) == name {
			return true
		}
	}
	return false
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
