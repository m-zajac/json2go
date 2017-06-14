package json2go

import "encoding/json"

// Parser parses successive json inputs and returns go representation as string
type Parser struct {
	rootNode *node
}

// NewParser creates new Parser
func NewParser() *Parser {
	rootNode := newNode("root")
	rootNode.root = true
	return &Parser{
		rootNode: rootNode,
	}
}

// FeedBytes consumes json input as bytes. If input is invalid, json unmarshalling error is returned
func (p *Parser) FeedBytes(input []byte) error {
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
func (p *Parser) FeedValue(input interface{}) {
	p.rootNode.grow(input)
}

// String returns string representation of go struct fitting parsed json values
func (p *Parser) String() (string, error) {
	return p.rootNode.repr()
}
