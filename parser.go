package jsontogo

type Parser struct {
	field field
}

func NewParser() *Parser {
	return &Parser{
		field: field{
			name: "root",
			root: true,
			t:    newInitType(),
		},
	}
}

func (p *Parser) Feed(input interface{}) {
	p.field.grow(input)
}

func (p *Parser) String() string {
	return p.field.repr()
}
