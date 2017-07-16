package json2go

const (
	nodeTypeInit          = nodeInitType(".")
	nodeTypeBool          = nodeBoolType("bool")
	nodeTypeInt           = nodeIntType("int")
	nodeTypeFloat         = nodeFloatType("float")
	nodeTypeString        = nodeStringType("string")
	nodeTypeUnknownObject = nodeUnknownObjectType("object?")
	nodeTypeObject        = nodeObjectType("object")
	nodeTypeInterface     = nodeInterfaceType("interface")
	nodeTypeExternal      = nodeExternalType("node")
)

type nodeType interface {
	id() string
	fit(interface{}) nodeType
	expands(nodeType) bool
}

func growType(t nodeType, v interface{}) nodeType {
	if v == nil {
		return t
	}

	new := t.fit(v)
	if t.id() != nodeTypeInit.id() && !new.expands(t) {
		return nodeTypeInterface
	}

	return new
}

// Type defs

type nodeInitType string

func (n nodeInitType) id() string {
	return string(n)
}

func (n nodeInitType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeInitType) fit(v interface{}) nodeType {
	return nodeTypeBool.fit(v)
}

type nodeBoolType string

func (n nodeBoolType) id() string {
	return string(n)
}

func (n nodeBoolType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeBoolType) fit(v interface{}) nodeType {
	switch v.(type) {
	case bool:
		return n
	}

	return nodeTypeInt.fit(v)
}

type nodeIntType string

func (n nodeIntType) id() string {
	return string(n)
}

func (n nodeIntType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeIntType) fit(v interface{}) nodeType {
	switch typedValue := v.(type) {
	case int, int8, int16, int32, int64:
		return n
	case float32:
		if typedValue == float32(int(typedValue)) {
			return n
		}
	case float64:
		if typedValue == float64(int(typedValue)) {
			return n
		}
	}

	return nodeTypeFloat.fit(v)
}

type nodeFloatType string

func (n nodeFloatType) id() string {
	return string(n)
}

func (n nodeFloatType) expands(n2 nodeType) bool {
	return n == n2 || n2.id() == nodeTypeInt.id()
}

func (n nodeFloatType) fit(v interface{}) nodeType {
	switch v.(type) {
	case float32, float64, int, int16, int32, int64:
		return n
	}

	return nodeTypeString.fit(v)
}

type nodeStringType string

func (n nodeStringType) id() string {
	return string(n)
}

func (n nodeStringType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeStringType) fit(v interface{}) nodeType {
	switch v.(type) {
	case string:
		return n
	}

	return nodeTypeUnknownObject.fit(v)
}

type nodeUnknownObjectType string

func (n nodeUnknownObjectType) id() string {
	return string(n)
}

func (n nodeUnknownObjectType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeUnknownObjectType) fit(v interface{}) nodeType {
	switch typedValue := v.(type) {
	case map[string]interface{}:
		if len(typedValue) == 0 {
			return n
		}
	}

	return nodeTypeObject.fit(v)
}

type nodeObjectType string

func (n nodeObjectType) id() string {
	return string(n)
}

func (n nodeObjectType) expands(n2 nodeType) bool {
	return n == n2 || n2 == nodeTypeUnknownObject
}

func (n nodeObjectType) fit(v interface{}) nodeType {
	switch v.(type) {
	case map[string]interface{}:
		return n
	}

	return nodeTypeInterface.fit(v)
}

type nodeInterfaceType string

func (n nodeInterfaceType) id() string {
	return string(n)
}

func (n nodeInterfaceType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeInterfaceType) fit(v interface{}) nodeType {
	return n
}

type nodeExternalType string

func (n nodeExternalType) id() string {
	return string(n)
}

func (n nodeExternalType) expands(n2 nodeType) bool {
	return n == n2
}

func (n nodeExternalType) fit(v interface{}) nodeType {
	return n
}
