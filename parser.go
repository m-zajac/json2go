package json2go

import (
	"encoding/json"
	"go/ast"
)

type options struct {
	extractCommonTypes           bool
	stringPointersWhenKeyMissing bool
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
	p.rootNode.sort()
	nodes := []*node{p.rootNode}
	if p.opts.extractCommonTypes {
		nodes = extractCommonSubtrees(p.rootNode)
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
