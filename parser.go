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

// Feed consumes json input as bytes. If input is invalid, json unmarshalling error is returned
func (p *Parser) Feed(input []byte) error {
	var v interface{}
	if err := json.Unmarshal(input, &v); err != nil {
		return err
	}

	p.field.grow(v)

	return nil
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
