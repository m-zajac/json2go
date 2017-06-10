package json2go

type fieldTypeID string

const (
	fieldTypeInit fieldTypeID = "init"

	fieldTypeBool          fieldTypeID = "bool"
	fieldTypeInt           fieldTypeID = "int"
	fieldTypeFloat         fieldTypeID = "float"
	fieldTypeString        fieldTypeID = "string"
	fieldTypeUnknownObject fieldTypeID = "object?"
	fieldTypeObject        fieldTypeID = "object"
	fieldTypeInterface     fieldTypeID = "interface"

	fieldTypeArrayUnknown   fieldTypeID = "[]?"
	fieldTypeArrayBool      fieldTypeID = "[]bool"
	fieldTypeArrayInt       fieldTypeID = "[]int"
	fieldTypeArrayFloat     fieldTypeID = "[]float"
	fieldTypeArrayString    fieldTypeID = "[]string"
	fieldTypeArrayObject    fieldTypeID = "[]object"
	fieldTypeArrayInterface fieldTypeID = "[]interface"
)

type fieldType struct {
	id        fieldTypeID
	fitFunc   func(fieldType, interface{}) fieldType
	arrayFunc func() fieldType
	reprFunc  func() string

	expandsTypes []fieldType
	grown        bool
}

func (k fieldType) grow(value interface{}) fieldType {
	if value == nil {
		return k
	}

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
	}

	return k.fitFunc(k, value)
}

func (k fieldType) expands(k2 fieldType) bool {
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
	return fieldType{
		id: fieldTypeInit,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			return newBoolType().fit(value)
		},
		arrayFunc: func() fieldType {
			return newUnknownArrayType()
		},
		reprFunc: func() string {
			return "interface{}"
		},
	}
}

func newBoolType() fieldType {
	return fieldType{
		id: fieldTypeBool,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case bool:
				return k
			}

			return newIntType().fit(value)
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
		id: fieldTypeInt,
		fitFunc: func(k fieldType, value interface{}) fieldType {
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
		id: fieldTypeFloat,
		expandsTypes: []fieldType{
			newIntType(),
		},
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case float32, float64:
				return k
			}

			return newStringType().fit(value)
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
		id: fieldTypeString,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case string:
				return k
			}

			return newUnknownObjectType().fit(value)
		},
		arrayFunc: func() fieldType {
			return newStringArrayType()
		},
		reprFunc: func() string {
			return "string"
		},
	}
}

func newUnknownObjectType() fieldType {
	return fieldType{
		id: fieldTypeUnknownObject,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch typedValue := value.(type) {
			case map[string]interface{}:
				if len(typedValue) == 0 {
					return k
				}
			}

			return newObjectType().fit(value)
		},
		arrayFunc: func() fieldType {
			return newInterfaceArrayType()
		},
		reprFunc: func() string {
			return "interface{}"
		},
	}
}

func newObjectType() fieldType {
	return fieldType{
		id: fieldTypeObject,
		expandsTypes: []fieldType{
			newUnknownObjectType(),
		},
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch value.(type) {
			case map[string]interface{}:
				return k
			}

			return newInterfaceType().fit(value)
		},
		arrayFunc: func() fieldType {
			return newObjectArrayType()
		},
		reprFunc: func() string {
			return "struct"
		},
	}
}

func newInterfaceType() fieldType {
	return fieldType{
		id: fieldTypeInterface,
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
		id: fieldTypeArrayUnknown,
		fitFunc: func(k fieldType, value interface{}) fieldType {
			switch typedValue := value.(type) {
			case []interface{}:
				return newArrayFieldTypeFromValues(typedValue)
			}

			return newInterfaceType()
		},
		reprFunc: func() string {
			return "[]interface{}"
		},
	}
}

func newBoolArrayType() fieldType {
	return fieldType{
		id:      fieldTypeArrayBool,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]bool"
		},
	}
}

func newIntArrayType() fieldType {
	return fieldType{
		id:      fieldTypeArrayInt,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]int"
		},
	}
}

func newFloatArrayType() fieldType {
	return fieldType{
		id:      fieldTypeArrayFloat,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]float64"
		},
	}
}

func newStringArrayType() fieldType {
	return fieldType{
		id:      fieldTypeArrayString,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]string"
		},
	}
}

func newObjectArrayType() fieldType {
	return fieldType{
		id:      fieldTypeArrayObject,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]struct"
		},
	}
}

func newInterfaceArrayType() fieldType {
	return fieldType{
		id: fieldTypeArrayInterface,
		expandsTypes: []fieldType{
			newBoolArrayType(),
			newIntArrayType(),
			newFloatArrayType(),
			newStringArrayType(),
			newObjectArrayType(),
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

func fitArray(k fieldType, value interface{}) fieldType {
	sliceValue, ok := value.([]interface{})
	if !ok {
		return newInterfaceType()
	}

	ak := newArrayFieldTypeFromValues(sliceValue)
	if k.expands(ak) {
		return k
	}
	if ak.expands(k) {
		return ak
	}

	return newInterfaceArrayType()
}
