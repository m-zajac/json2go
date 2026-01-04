package json2go

import (
	"encoding/json"
	"go/ast"
)

type options struct {
	// extractAllTypes toggles extracting all nested object types as separate types.
	extractAllTypes bool
	// extractCommonTypes toggles extracting common JSON nodes as separate types.
	extractCommonTypes bool
	// extractSimilarityThreshold is the minimum similarity score (0.0 to 1.0) required
	// to consider two types as similar enough to be merged.
	extractSimilarityThreshold float64
	// extractMinSubsetSize is the minimum number of fields a shared subset must have
	// to be considered for extraction as a separate type.
	extractMinSubsetSize int
	// extractMinSubsetOccurrences is the minimum number of times a shared subset
	// must occur across different objects to be extracted.
	extractMinSubsetOccurrences int
	// extractMinAddedFields is the minimum number of unique fields a new type must have
	// (excluding fields from other embedded types) to justify its extraction.
	extractMinAddedFields int
	// stringPointersWhenKeyMissing toggles whether missing string keys in some documents
	// should result in a pointer (*string) instead of a regular string.
	stringPointersWhenKeyMissing bool
	// skipEmptyKeys toggles skipping keys in the input that contained only null values.
	skipEmptyKeys bool
	// makeMaps defines if the parser should attempt to use maps instead of structs
	// when objects have similar values and can be represented as map[string]T.
	makeMaps bool
	// makeMapsWhenMinAttributes defines the minimum number of attributes an object
	// must have to be considered for conversion to a map.
	makeMapsWhenMinAttributes uint
	// timeAsStr toggles whether to treat valid time strings as time.Time or just as strings.
	timeAsStr bool
}

// JSONParserOpt is a type for setting parser options.
type JSONParserOpt func(*options)

// OptExtractAllTypes toggles extracting all nested object types as separate types.
func OptExtractAllTypes(v bool) JSONParserOpt {
	return func(o *options) {
		o.extractAllTypes = v
	}
}

// OptExtractCommonTypes toggles extracting common JSON nodes as separate types.
func OptExtractCommonTypes(v bool) JSONParserOpt {
	return func(o *options) {
		o.extractCommonTypes = v
	}
}

// OptExtractHeuristics sets thresholds for extracting similar types and shared subsets.
// similarity is the minimum score (0.0 to 1.0) for merging types.
// minSize is the minimum number of fields in a shared subset.
// minOccurrences is the minimum frequency of a subset to be extracted.
// minAddedFields is the minimum unique fields required for a new type.
func OptExtractHeuristics(similarity float64, minSize, minOccurrences, minAddedFields int) JSONParserOpt {
	return func(o *options) {
		o.extractSimilarityThreshold = similarity
		o.extractMinSubsetSize = minSize
		o.extractMinSubsetOccurrences = minOccurrences
		o.extractMinAddedFields = minAddedFields
	}
}

// OptStringPointersWhenKeyMissing toggles whether missing string keys in some documents
// should result in a pointer (*string) instead of a regular string.
func OptStringPointersWhenKeyMissing(v bool) JSONParserOpt {
	return func(o *options) {
		o.stringPointersWhenKeyMissing = v
	}
}

// OptSkipEmptyKeys toggles skipping keys in the input that contained only null values.
func OptSkipEmptyKeys(v bool) JSONParserOpt {
	return func(o *options) {
		o.skipEmptyKeys = v
	}
}

// OptMakeMaps defines if the parser should attempt to use maps instead of structs
// when objects have similar values and can be represented as map[string]T.
// minAttributes defines the minimum number of attributes an object must have
// to be considered for conversion to a map.
func OptMakeMaps(v bool, minAttributes uint) JSONParserOpt {
	return func(o *options) {
		o.makeMaps = v
		o.makeMapsWhenMinAttributes = minAttributes
	}
}

// OptTimeAsString toggles whether to treat valid time strings as time.Time or just as strings.
func OptTimeAsString(v bool) JSONParserOpt {
	return func(o *options) {
		o.timeAsStr = v
	}
}

// JSONParser parses successive json inputs and returns go representation as string
type JSONParser struct {
	rootNode *node
	opts     options
}

// NewJSONParser creates new json Parser
func NewJSONParser(rootTypeName string, opts ...JSONParserOpt) *JSONParser {
	rootNode := newNode(rootTypeName)
	rootNode.root = true
	p := JSONParser{
		rootNode: rootNode,
		opts:     defaultOptions(),
	}
	for _, o := range opts {
		o(&p.opts)
	}

	return &p
}

func defaultOptions() options {
	return options{
		extractAllTypes:              DefaultExtractAllTypes,
		extractCommonTypes:           DefaultExtractCommonTypes,
		extractSimilarityThreshold:   DefaultSimilarityThreshold,
		extractMinSubsetSize:         DefaultMinSubsetSize,
		extractMinSubsetOccurrences:  DefaultMinSubsetOccurrences,
		extractMinAddedFields:        DefaultMinAddedFields,
		stringPointersWhenKeyMissing: DefaultStringPointersWhenKeyMissing,
		skipEmptyKeys:                DefaultSkipEmptyKeys,
		makeMaps:                     DefaultMakeMaps,
		makeMapsWhenMinAttributes:    uint(DefaultMakeMapsWhenMinAttributes),
		timeAsStr:                    DefaultTimeAsStr,
	}
}

// FeedBytes consumes json input as bytes. If input is invalid, json unmarshalling error is returned
func (p *JSONParser) FeedBytes(input []byte) error {
	var v interface{}
	if err := json.Unmarshal(input, &v); err != nil {
		return err
	}

	p.FeedValue(v)

	return nil
}

// FeedValue consumes one of:
//
//   - simple type (int, float, string, etc.)
//   - []interface{} - each value must meet these requirements
//   - map[string]interface{}  - each value must meet these requirements
//
// json.Unmarshal to empty interface value provides perfect input (see example)
func (p *JSONParser) FeedValue(input interface{}) {
	p.rootNode.grow(input)
}

// String returns string representation of go struct fitting parsed json values
func (p *JSONParser) String() string {
	root := p.rootNode.clone()

	root.sort()

	if p.opts.skipEmptyKeys {
		p.stripEmptyKeys(root)
	}
	if p.opts.makeMaps {
		convertViableObjectsToMaps(root, p.opts.makeMapsWhenMinAttributes)
	}

	nodes := []*node{root}
	if p.opts.extractCommonTypes {
		nodes = extractCommonSubtrees(root, p.opts)
	}
	if p.opts.extractAllTypes {
		// Extract any remaining inline types that weren't extracted by extractCommonTypes
		nodes = extractRemainingNestedTypes(nodes, p.opts)
	}

	// Sort nodes in topological order
	nodes = sortNodesBFS(root, nodes)

	return astPrintDecls(
		astMakeDecls(nodes, p.opts),
	)
}

// ASTDecls returns ast type declarations
func (p *JSONParser) ASTDecls() []ast.Decl {
	p.rootNode.sort()
	return astMakeDecls(
		[]*node{p.rootNode},
		p.opts,
	)
}

func (p *JSONParser) stripEmptyKeys(n *node) {
	if len(n.children) == 0 {
		return
	}

	newChildren := make([]*node, 0, len(n.children))
	for i, c := range n.children {
		if c.t.id() != nodeTypeInit.id() {
			p.stripEmptyKeys(c)
			newChildren = append(newChildren, n.children[i])
		}
	}
	n.children = newChildren
}

// sortNodesBFS performs breadth-first topological sort on nodes,
// ordering them by dependency level with ordering by appearance at each level.
// The root node comes first, followed by its direct dependencies (in order of appearance),
// then their dependencies, and so on.
func sortNodesBFS(root *node, nodes []*node) []*node {
	if len(nodes) <= 1 {
		return nodes
	}

	// Build a map of node name -> node for quick lookup
	nodeMap := make(map[string]*node)
	for _, n := range nodes {
		nodeMap[n.name] = n
	}

	// Track visited nodes and result order
	visited := make(map[string]bool)
	var result []*node

	// BFS queue - we'll process level by level
	var queue []*node
	queue = append(queue, root)
	visited[root.name] = true

	for len(queue) > 0 {
		// Process current level
		levelSize := len(queue)
		var currentLevel []*node

		// Collect all dependencies from current level in order of appearance
		var deps []string
		seenDeps := make(map[string]bool)
		for i := range levelSize {
			node := queue[i]
			currentLevel = append(currentLevel, node)

			// Collect all dependencies of this node in order
			nodeDeps := collectDependenciesOrdered(node)
			for _, depName := range nodeDeps {
				if !visited[depName] && !seenDeps[depName] {
					if _, exists := nodeMap[depName]; exists {
						deps = append(deps, depName)
						seenDeps[depName] = true
					}
				}
			}
		}

		// Remove current level from queue
		queue = queue[levelSize:]

		// Add current level to result
		result = append(result, currentLevel...)

		// Add dependencies to queue in order of appearance
		for _, depName := range deps {
			if !visited[depName] {
				visited[depName] = true
				queue = append(queue, nodeMap[depName])
			}
		}
	}

	// Handle any orphan nodes (not referenced by anyone)
	var orphans []*node
	for _, n := range nodes {
		if !visited[n.name] {
			orphans = append(orphans, n)
		}
	}

	// Sort orphans alphabetically by name
	sortNodesByName(orphans)
	result = append(result, orphans...)

	return result
}

// collectDependenciesOrdered recursively finds all extracted type names referenced by a node,
// preserving the order of appearance in the node's children.
func collectDependenciesOrdered(n *node) []string {
	if n == nil {
		return nil
	}

	var deps []string
	seen := make(map[string]bool)

	var collect func(*node)
	collect = func(node *node) {
		if node == nil {
			return
		}
		for _, child := range node.children {
			if child.t.id() == nodeTypeExtracted.id() && child.extractedTypeName != "" {
				if !seen[child.extractedTypeName] {
					deps = append(deps, child.extractedTypeName)
					seen[child.extractedTypeName] = true
				}
			}
			// Recursively check nested children
			collect(child)
		}
	}

	collect(n)
	return deps
}

// sortNodesByName sorts nodes by their name field alphabetically
func sortNodesByName(nodes []*node) {
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if nodes[i].name > nodes[j].name {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}
