package json2go

const (
	rootSchema   = "http://json-schema.org/draft-07/schema"
	rootSchemaID = "http://example.com/example.json"
)

type schemaItem struct {
	Schema     string       `json:"$schema,omitempty"`
	ID         string       `json:"$id"`
	Type       string       `json:"type"`
	Properties []schemaItem `json:"properties,omitempty"`
	Required   []string     `json:"required,omitempty"`
}

func newSchemaItemFromNode(n node) schemaItem {
	item := schemaItem{
		Type: schemaType(n.t),
	}
	if n.root {
		item.Schema = rootSchema
		item.ID = rootSchemaID
	}

	return item
}

func schemaType(t nodeType) string {
	switch t.id() {
	case nodeTypeInit.id():
		fallthrough
	case nodeTypeUnknownObject.id():
		fallthrough
	case nodeTypeInterface.id():
		fallthrough
	case nodeTypeExternal.id():
		fallthrough
	case nodeTypeObject.id():
		return "object"
	case nodeTypeBool.id():
		return "boolean"
	case nodeTypeInt.id():
		return "integer"
	case nodeTypeFloat.id():
		return "number"
	case nodeTypeString.id():
		return "string"
	default:
		return "object"
	}
}
