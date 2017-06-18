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

	nodeTypeArrayUnknown   nodeTypeID = "[]?"
	nodeTypeArrayBool      nodeTypeID = "[]bool"
	nodeTypeArrayInt       nodeTypeID = "[]int"
	nodeTypeArrayFloat     nodeTypeID = "[]float"
	nodeTypeArrayString    nodeTypeID = "[]string"
	nodeTypeArrayObject    nodeTypeID = "[]object"
	nodeTypeArrayInterface nodeTypeID = "[]interface"

	nodeTypeExternalNode nodeTypeID = "node" // this node type points to other tree root
)

type nodeType struct {
	id        nodeTypeID
	fitFunc   func(nodeType, interface{}) nodeType
	arrayFunc func() nodeType

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
	switch typedValue := value.(type) {
	case []interface{}:
		if k.arrayFunc != nil { // k is base type
			return newArrayNodeTypeFromValues(typedValue)
		}
	}

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

func (k nodeType) arrayType() nodeType {
	if k.arrayFunc == nil {
		return newUnknownArrayType()
	}
	return k.arrayFunc()
}

func newInitType() nodeType {
	return nodeType{
		id: nodeTypeInit,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			return newBoolType().fit(value)
		},
		arrayFunc: func() nodeType {
			return newUnknownArrayType()
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
		arrayFunc: func() nodeType {
			return newBoolArrayType()
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
		arrayFunc: func() nodeType {
			return newIntArrayType()
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
			case float32, float64:
				return k
			}

			return newStringType().fit(value)
		},
		arrayFunc: func() nodeType {
			return newFloatArrayType()
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
		arrayFunc: func() nodeType {
			return newStringArrayType()
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
		arrayFunc: func() nodeType {
			return newInterfaceArrayType()
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
		arrayFunc: func() nodeType {
			return newObjectArrayType()
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

func newUnknownArrayType() nodeType {
	return nodeType{
		id: nodeTypeArrayUnknown,
		fitFunc: func(k nodeType, value interface{}) nodeType {
			switch typedValue := value.(type) {
			case []interface{}:
				return newArrayNodeTypeFromValues(typedValue)
			}

			return newInterfaceType()
		},
	}
}

func newBoolArrayType() nodeType {
	return nodeType{
		id:      nodeTypeArrayBool,
		fitFunc: fitArray,
	}
}

func newIntArrayType() nodeType {
	return nodeType{
		id:      nodeTypeArrayInt,
		fitFunc: fitArray,
	}
}

func newFloatArrayType() nodeType {
	return nodeType{
		id:      nodeTypeArrayFloat,
		fitFunc: fitArray,
	}
}

func newStringArrayType() nodeType {
	return nodeType{
		id:      nodeTypeArrayString,
		fitFunc: fitArray,
	}
}

func newObjectArrayType() nodeType {
	return nodeType{
		id:      nodeTypeArrayObject,
		fitFunc: fitArray,
	}
}

func newInterfaceArrayType() nodeType {
	return nodeType{
		id: nodeTypeArrayInterface,
		expandsTypes: []nodeType{
			newBoolArrayType(),
			newIntArrayType(),
			newFloatArrayType(),
			newStringArrayType(),
			newObjectArrayType(),
		},
		fitFunc: fitArray,
	}
}

func newExternalObjectType() nodeType {
	return nodeType{
		id: nodeTypeExternalNode,
	}
}

func newArrayNodeTypeFromValues(values []interface{}) nodeType {
	if values == nil || len(values) == 0 {
		return newUnknownArrayType()
	}

	var valuesTypes []nodeType
	for _, v := range values {
		valuesTypes = append(valuesTypes, newInitType().fit(v))
	}

	selectedType := valuesTypes[0]

loop:
	for _, k := range valuesTypes {
		if k.id == selectedType.id {
			continue
		}

		if k.expands(selectedType) {
			selectedType = k
			continue loop
		}
		if selectedType.expands(k) {
			continue loop
		}

		return newInterfaceArrayType()
	}

	return selectedType.arrayType()
}

func fitArray(k nodeType, value interface{}) nodeType {
	sliceValue, ok := value.([]interface{})
	if !ok {
		return newInterfaceType()
	}

	ak := newArrayNodeTypeFromValues(sliceValue)
	if k.expands(ak) {
		return k
	}
	if ak.expands(k) {
		return ak
	}

	return newInterfaceArrayType()
}
