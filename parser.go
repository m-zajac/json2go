package json2go

import "encoding/json"

// Parser parses successive json inputs and returns go representation as string
type Parser struct {
	field *field
}

// NewParser creates new Parser
func NewParser() *Parser {
	rootField := newField("root")
	rootField.root = true
	return &Parser{
		field: rootField,
	}
}

// FeedBytes consumes json input as bytes. If input is invalid, json unmarshalling error is returned
func (p *Parser) FeedBytes(input []byte) error {
	var v interface{}
	if err := json.Unmarshal(input, &v); err != nil {
		return err
	}

	p.field.grow(v)

	return nil
}

// FeedObject consumes any value coming from interface value.
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
// 	parser.FeedObject(v)
func (p *Parser) FeedObject(input interface{}) {
	p.field.grow(input)
}

// Result returns string representation of go struct fitting parsed json values
func (p *Parser) Result() (string, error) {
	return p.field.repr()
}

// Strings acts as Result (fmt.Stringer)
func (p *Parser) String() string {
	s, _ := p.field.repr()
	return s
}
