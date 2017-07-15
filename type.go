package json2go

type nodeTypeID string

const (
	nodeTypeInit nodeTypeID = "init"

	nodeTypeBool          nodeTypeID = "bool"
	nodeTypeInt           nodeTypeID = "int"
	nodeTypeFloat         nodeTypeID = "float"
	nodeTypeString        nodeTypeID = "string"
	nodeTypeUnknownObject nodeTypeID = "object?"
	nodeTypeObject        nodeTypeID = "object"
	nodeTypeInterface     nodeTypeID = "interface"

	nodeTypeExternalNode nodeTypeID = "node" // this node type points to other tree root
)

type nodeType struct {
	id      nodeTypeID
	fitFunc func(nodeType, interface{}) nodeType

	expandsTypes []nodeType
}

func (k nodeType) grow(value interface{}) nodeType {
	if value == nil {
		return k
	}

	new := k.fit(value)
	if k.id != nodeTypeInit && !new.expands(k) {
		return newInterfaceType()
	}

	return new
}

func (k nodeType) fit(value interface{}) nodeType {
	return k.fitFunc(k, value)
}

func (k nodeType) expands(k2 nodeType) bool {
	if k.id == k2.id {
		return true
	}

	for _, smallerType := range k.expandsTypes {
		if k2.id == smallerType.id {
			return smallerType.expands(k2)
		}
	}

	return false
}

func newInitType() nodeType {
	return nodeType{
		id: nodeTypeInit,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			return newBoolType().fit(value)
		},
	}
}

func newBoolType() nodeType {
	return nodeType{
		id: nodeTypeBool,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch value.(type) {
			case bool:
				return k
			}

			return newIntType().fit(value)
		},
	}
}

func newIntType() nodeType {
	return nodeType{
		id: nodeTypeInt,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch typedValue := value.(type) {
			case int, int8, int16, int32, int64:
				return k
			case float32:
				if typedValue == float32(int(typedValue)) {
					return k
				}
			case float64:
				if typedValue == float64(int(typedValue)) {
					return k
				}
			}
			k = newFloatType()
			return k.fit(value)
		},
	}
}

func newFloatType() nodeType {
	return nodeType{
		id: nodeTypeFloat,
		expandsTypes: []nodeType{
			newIntType(),
		},
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch value.(type) {
			case float32, float64, int, int16, int32, int64:
				return k
			}

			return newStringType().fit(value)
		},
	}
}

func newStringType() nodeType {
	return nodeType{
		id: nodeTypeString,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch value.(type) {
			case string:
				return k
			}

			return newUnknownObjectType().fit(value)
		},
	}
}

func newUnknownObjectType() nodeType {
	return nodeType{
		id: nodeTypeUnknownObject,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch typedValue := value.(type) {
			case map[string]interface{}:
				if len(typedValue) == 0 {
					return k
				}
			}

			return newObjectType().fit(value)
		},
	}
}

func newObjectType() nodeType {
	return nodeType{
		id: nodeTypeObject,
		expandsTypes: []nodeType{
			newUnknownObjectType(),
		},
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch value.(type) {
			case map[string]interface{}:
				return k
			}

			return newInterfaceType().fit(value)
		},
	}
}

func newInterfaceType() nodeType {
	return nodeType{
		id: nodeTypeInterface,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			return k
		},
	}
}

func newExternalObjectType() nodeType {
	return nodeType{
		id: nodeTypeExternalNode,
	}
}
