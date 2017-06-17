package json2go

import (
	"encoding/json"
	"go/ast"
)

// JSONParser parses successive json inputs and returns go representation as string
type JSONParser struct {
	rootNode *node
}

// NewJSONParser creates new json Parser
func NewJSONParser() *JSONParser {
	rootNode := newNode("object")
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

	p.rootNode.grow(v)

	return nil
}

// FeedValue consumes any value coming from interface value.
// If input can be one of:
//
//	* simple type (int, float, string, etc.)
//	* []interface{} - each value must meet these requirements
//	* map[string]interface{}  - each value must meet these requirements
//
// json.Unmarshal to empty interface value provides perfect input.
//
// Example:
// 	var v interface{}
// 	if err := json.Unmarshal(input, &v); err != nil {
// 		return err
// 	}
// 	parser.FeedValue(v)
func (p *JSONParser) FeedValue(input interface{}) {
	p.rootNode.grow(input)
}

// String returns string representation of go struct fitting parsed json values
func (p *JSONParser) String() string {
	return astPrintDecls(astMakeDecls(p.rootNode))
}

// ASTDecls returns ast type declarations
func (p *JSONParser) ASTDecls() []ast.Decl {
	return astMakeDecls(p.rootNode)
}
