package json2go

import (
	"encoding/json"
	"go/ast"
)

// JSONParser parses successive json inputs and returns go representation as string
type JSONParser struct {
	rootNode *node

	ExtractCommonStructs bool
}

// NewJSONParser creates new json Parser
func NewJSONParser(rootTypeName string) *JSONParser {
	rootNode := newNode(rootTypeName)
	rootNode.root = true
	return &JSONParser{
		rootNode: rootNode,
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
//	* simple type (int, float, string, etc.)
//	* []interface{} - each value must meet these requirements
//	* map[string]interface{}  - each value must meet these requirements
//
// json.Unmarshal to empty interface value provides perfect input (see example)
func (p *JSONParser) FeedValue(input interface{}) {
	if slice, ok := input.([]interface{}); ok {
		for _, v := range slice {
			p.rootNode.grow(v)
		}
		return
	}

	p.rootNode.grow(input)
}

// String returns string representation of go struct fitting parsed json values
func (p *JSONParser) String() string {
	p.rootNode.sort()
	nodes := []*node{p.rootNode}
	if p.ExtractCommonStructs {
		nodes = extractCommonSubtrees(p.rootNode)
	}
	return astPrintDecls(astMakeDecls(nodes))
}

// ASTDecls returns ast type declarations
func (p *JSONParser) ASTDecls() []ast.Decl {
	p.rootNode.sort()
	return astMakeDecls([]*node{p.rootNode})
}
