package json2go

import (
	"encoding/json"
	"go/ast"
)

type options struct {
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
