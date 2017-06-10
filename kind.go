package jsontogo

type kindName string

const (
	kindNameBool      kindName = "bool"
	kindNameInt       kindName = "int"
	kindNameFloat     kindName = "float"
	kindNameString    kindName = "string"
	kindNameInterface kindName = "interface"

	kindNameArrayUnknown   kindName = "aunknown"
	kindNameArrayBool      kindName = "abool"
	kindNameArrayInt       kindName = "aint"
	kindNameArrayFloat     kindName = "afloat"
	kindNameArrayString    kindName = "astring"
	kindNameArrayInterface kindName = "ainterface"
)

type kind struct {
	name      kindName
	fitFunc   func(kind, interface{}) kind
	arrayFunc func() kind
	reprFunc  func() string

	expandsKinds []kind
	grown        bool
}

func (k kind) grow(value interface{}) kind {
	new := k.fit(value)
	if k.grown && !new.expands(k) {
		return newInterfaceKind()
	}
	new.grown = true

	return new
}

func (k kind) fit(value interface{}) kind {
	switch typedValue := value.(type) {
	case []interface{}:
		if k.arrayFunc != nil { // k is base type
			return newArrayKindFromValues(typedValue)
		}
		res := k.fitFunc(k, value)
		return res
	default:
		res := k.fitFunc(k, value)
		return res
	}
}

func (k kind) expands(k2 kind) bool {
	if k.name == k2.name {
		return true
	}

	for _, smallerKind := range k.expandsKinds {
		if k2.name == smallerKind.name {
			return smallerKind.expands(k2)
		}
	}

	return false
}

func (k kind) arrayKind() kind {
	if k.arrayFunc == nil {
		return newUnknownArrayKind()
	}
	return k.arrayFunc()
}

func (k kind) repr() string {
	return k.reprFunc()
}

func newStartingKind() kind {
	return newBoolKind()
}

func newBoolKind() kind {
	return kind{
		name: kindNameBool,
		fitFunc: func(k kind, value interface{}) kind {
			switch value.(type) {
			case bool:
				return k
			default:
				k = newIntKind()
				return k.fit(value)
			}
		},
		arrayFunc: func() kind {
			return newBoolArrayKind()
		},
		reprFunc: func() string {
			return "bool"
		},
	}
}

func newIntKind() kind {
	return kind{
		name: kindNameInt,
		fitFunc: func(k kind, value interface{}) kind {
			switch value.(type) {
			case int, int8, int16, int32, int64:
				return k
			default:
				k = newFloatKind()
				return k.fit(value)
			}
		},
		arrayFunc: func() kind {
			return newIntArrayKind()
		},
		reprFunc: func() string {
			return "int"
		},
	}
}

func newFloatKind() kind {
	return kind{
		name: kindNameFloat,
		expandsKinds: []kind{
			newIntKind(),
		},
		fitFunc: func(k kind, value interface{}) kind {
			switch value.(type) {
			case float32, float64:
				return k
			default:
				k = newStringKind()
				return k.fit(value)
			}
		},
		arrayFunc: func() kind {
			return newFloatArrayKind()
		},
		reprFunc: func() string {
			return "float64"
		},
	}
}

func newStringKind() kind {
	return kind{
		name: kindNameString,
		fitFunc: func(k kind, value interface{}) kind {
			switch value.(type) {
			case string:
				return k
			default:
				k = newInterfaceKind()
				return k.fit(value)
			}
		},
		arrayFunc: func() kind {
			return newStringArrayKind()
		},
		reprFunc: func() string {
			return "string"
		},
	}
}

func newInterfaceKind() kind {
	return kind{
		name: kindNameInterface,
		fitFunc: func(k kind, value interface{}) kind {
			return k
		},
		reprFunc: func() string {
			return "interface{}"
		},
	}
}

func newUnknownArrayKind() kind {
	return kind{
		name: kindNameArrayUnknown,
		fitFunc: func(k kind, value interface{}) kind {
			switch typedValue := value.(type) {
			case []interface{}:
				return newArrayKindFromValues(typedValue)
			default:
				return newInterfaceKind()
			}
		},
		reprFunc: func() string {
			return "[]interface{}"
		},
	}
}

func newBoolArrayKind() kind {
	return kind{
		name:    kindNameArrayBool,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]bool"
		},
	}
}

func newIntArrayKind() kind {
	return kind{
		name:    kindNameArrayInt,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]int"
		},
	}
}

func newFloatArrayKind() kind {
	return kind{
		name:    kindNameArrayFloat,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]float64"
		},
	}
}

func newStringArrayKind() kind {
	return kind{
		name:    kindNameArrayString,
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]string"
		},
	}
}

func newInterfaceArrayKind() kind {
	return kind{
		name: kindNameArrayInterface,
		expandsKinds: []kind{
			newBoolArrayKind(),
			newIntArrayKind(),
			newFloatArrayKind(),
			newStringArrayKind(),
		},
		fitFunc: fitArray,
		reprFunc: func() string {
			return "[]interface{}"
		},
	}
}

func newArrayKindFromValues(values []interface{}) kind {
	if values == nil || len(values) == 0 {
		return newUnknownArrayKind()
	}

	var valuesKinds []kind
	for _, v := range values {
		valuesKinds = append(valuesKinds, newStartingKind().fit(v))
	}

	selectedKind := valuesKinds[0]

loop:
	for _, k := range valuesKinds {
		if k.name == selectedKind.name {
			continue
		}

		if k.expands(selectedKind) {
			selectedKind = k
			continue loop
		}
		if selectedKind.expands(k) {
			continue loop
		}

		return newInterfaceArrayKind()
	}

	return selectedKind.arrayKind()
}

func fitArray(k kind, value interface{}) kind {
	switch typedValue := value.(type) {
	case []interface{}:
		ak := newArrayKindFromValues(typedValue)
		if k.expands(ak) {
			return k
		}
		if ak.expands(k) {
			return ak
		}

		return newInterfaceArrayKind()
	}

	return newInterfaceKind()
}
