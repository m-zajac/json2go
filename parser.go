package json2go

import (
	"encoding/json"
	"go/ast"
)

type options struct {
	extractCommonTypes           bool
	stringPointersWhenKeyMissing bool
	skipEmptyKeys                bool
	makeMaps                     bool
	makeMapsWhenMinAttributes    uint
	timeAsStr                    bool
}

// JSONParserOpt is a type for setting parser options.
type JSONParserOpt func(*options)

// OptExtractCommonTypes toggles extracting common json nodes as separate types.
func OptExtractCommonTypes(v bool) JSONParserOpt {
	return func(o *options) {
		o.extractCommonTypes = v
	}
}

// OptStringPointersWhenKeyMissing toggles wether missing string key in one of documents should result in pointer string.
func OptStringPointersWhenKeyMissing(v bool) JSONParserOpt {
	return func(o *options) {
		o.stringPointersWhenKeyMissing = v
	}
}

// OptSkipEmptyKeys toggles skipping keys in input that were only nulls.
func OptSkipEmptyKeys(v bool) JSONParserOpt {
	return func(o *options) {
		o.skipEmptyKeys = v
	}
}

// OptMakeMaps defines if parser should try to use maps instead of structs when possible.
// minAttributes defines minimum number of attributes in object to try converting it to a map.
func OptMakeMaps(v bool, minAttributes uint) JSONParserOpt {
	return func(o *options) {
		o.makeMaps = v
		o.makeMapsWhenMinAttributes = minAttributes
	}
}

// OptTimeAsString toggles using time.Time for valid time strings or just a string.
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
		opts:     options{},
	}
	for _, o := range opts {
		o(&p.opts)
	}

	return &p
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
//	* simple type (int, float, string, etc.)
//	* []interface{} - each value must meet these requirements
//	* map[string]interface{}  - each value must meet these requirements
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
		nodes = extractCommonSubtrees(root)
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

func (p *JSONParser) ASTDeclsWithOpt() []ast.Decl {
	root := p.rootNode.clone()
	root.sort()
	if p.opts.skipEmptyKeys {
		p.stripEmptyKeys(root)
	}
	if p.opts.makeMaps {
		convertViableObjectsToMaps(root, p.opts.makeMapsWhenMinAttributes)
	}
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
