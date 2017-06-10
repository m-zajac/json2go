package jsontogo

type fieldTypeID string

const (
	fieldTypeBool      fieldTypeID = "bool"
	fieldTypeInt       fieldTypeID = "int"
	fieldTypeFloat     fieldTypeID = "float"
	fieldTypeString    fieldTypeID = "string"
	fieldTypeInterface fieldTypeID = "interface"

	fieldTypeArrayUnknown   fieldTypeID = "[]unknown"
	fieldTypeArrayBool      fieldTypeID = "[]bool"
	fieldTypeArrayInt       fieldTypeID = "[]int"
	fieldTypeArrayFloat     fieldTypeID = "[]float"
	fieldTypeArrayString    fieldTypeID = "[]string"
	fieldTypeArrayInterface fieldTypeID = "[]interface"

	fieldTypeObject fieldTypeID = "object"
)

type fieldType struct {
	name      fieldTypeID
	fitFunc   func(fieldType, interface{}) fieldType
	arrayFunc func() fieldType
	reprFunc  func() string

	expandsTypes []fieldType
	grown        bool
}

func (k fieldType) grow(value interface{}) fieldType {
	new := k.fit(value)
	if k.grown && !new.expands(k) {
		return newInterfaceType()
	}
	new.grown = true

	return new
}

func (k fieldType) fit(value interface{}) fieldType {
	switch typedValue := value.(type) {
	case []interface{}:
		if k.arrayFunc != nil { // k is base type
			return newArrayFieldTypeFromValues(typedValue)
		}
		res := k.fitFunc(k, value)
		return res
	default:
		res := k.fitFunc(k, value)
		return res
	}
}

func (k fieldType) expands(k2 fieldType) bool {
	if k.name == k2.name {
		return true
	}

	for _, smallerType := range k.expandsTypes {
		if k2.name == smallerType.name {
			return smallerType.expands(k2)
		}
	}

	return false
}

func (k fieldType) arrayType() fieldType {
	if k.arrayFunc == nil {
		return newUnknownArrayType()
	}
	return k.arrayFunc()
}

func (k fieldType) repr() string {
	return k.reprFunc()
}

func newInitType() fieldType {
	return newBoolType()
}

func newBoolType() fieldType {
	return fieldType{
		name: fieldTypeBool,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case bool:
				return k
			default:
				k = newIntType()
				return k.fit(value)
			}
		},
		arrayFunc: func() fieldType {
			return newBoolArrayType()
		},
		reprFunc: func() string {
			return "bool"
		},
	}
}

func newIntType() fieldType {
	return fieldType{
		name: fieldTypeInt,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case int, int8, int16, int32, int64:
				return k
			default:
				k = newFloatType()
				return k.fit(value)
			}
		},
		arrayFunc: func() fieldType {
			return newIntArrayType()
		},
		reprFunc: func() string {
			return "int"
		},
	}
}

func newFloatType() fieldType {
	return fieldType{
		name: fieldTypeFloat,
		expandsTypes: []fieldType{
			newIntType(),
		},
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case float32, float64:
				return k
			default:
				k = newStringType()
				return k.fit(value)
			}
		},
		arrayFunc: func() fieldType {
			return newFloatArrayType()
		},
		reprFunc: func() string {
			return "float64"
		},
	}
}

func newStringType() fieldType {
	return fieldType{
		name: fieldTypeString,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case string:
				return k
			default:
				k = newObjectType()
				return k.fit(value)
			}
		},
		arrayFunc: func() fieldType {
			return newStringArrayType()
		},
		reprFunc: func() string {
			return "string"
		},
	}
}

func newObjectType() fieldType {
	return fieldType{
		name: fieldTypeObject,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case map[string]interface{}:
				return k
			default:
				k = newInterfaceType()
				return k.fit(value)
			}
		},
		reprFunc: func() string {
			return "struct"
		},
	}
}

func newInterfaceType() fieldType {
	return fieldType{
		name: fieldTypeInterface,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			return k
		},
		reprFunc: func() string {
			return "interface{}"
		},
	}
}

func newUnknownArrayType() fieldType {
	return fieldType{
		name: fieldTypeArrayUnknown,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch typedValue := value.(type) {
			case []interface{}:
				return newArrayFieldTypeFromValues(typedValue)
			default:
				return newInterfaceType()
			}
		},
		reprFunc: func() string {
			return "[]interface{}"
		},
	}
}

func newBoolArrayType() fieldType {
	return fieldType{
		name:    fieldTypeArrayBool,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]bool"
		},
	}
}

func newIntArrayType() fieldType {
	return fieldType{
		name:    fieldTypeArrayInt,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]int"
		},
	}
}

func newFloatArrayType() fieldType {
	return fieldType{
		name:    fieldTypeArrayFloat,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]float64"
		},
	}
}

func newStringArrayType() fieldType {
	return fieldType{
		name:    fieldTypeArrayString,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]string"
		},
	}
}

func newInterfaceArrayType() fieldType {
	return fieldType{
		name: fieldTypeArrayInterface,
		expandsTypes: []fieldType{
			newBoolArrayType(),
			newIntArrayType(),
			newFloatArrayType(),
			newStringArrayType(),
		},
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]interface{}"
		},
	}
}

func newArrayFieldTypeFromValues(values []interface{}) fieldType {
	if values == nil || len(values) == 0 {
		return newUnknownArrayType()
	}

	var valuesTypes []fieldType
	for _, v := range values {
		valuesTypes = append(valuesTypes, newInitType().fit(v))
	}

	selectedType := valuesTypes[0]

loop:
	for _, k := range valuesTypes {
		if k.name == selectedType.name {
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

func fitArray(k fieldType, value interface{}) fieldType {
	switch typedValue := value.(type) {
	case []interface{}:
		ak := newArrayFieldTypeFromValues(typedValue)
		if k.expands(ak) {
			return k
		}
		if ak.expands(k) {
			return ak
		}

		return newInterfaceArrayType()
	}

	return newInterfaceType()
}
