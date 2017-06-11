package json2go

import (
	"encoding/json"
	"fmt"
	"io"
)

// Decode reads json data from input reader and writes go type definition to given writer
func Decode(r io.Reader, w io.Writer) error {
	var data interface{}

	jd := json.NewDecoder(r)
	if err := jd.Decode(&data); err != nil {
		return fmt.Errorf("json decoding error: %v", err)
	}

	rootField := newField("root")
	rootField.root = true

	switch typedData := data.(type) {
	case []interface{}:
		for _, v := range typedData {
			rootField.grow(v)
		}
	default:
		rootField.grow(data)
	}

	repr, err := rootField.repr()
	if err != nil {
		return fmt.Errorf("generating type error: %v", err)
	}

	w.Write([]byte("\n"))
	w.Write([]byte(repr))
	w.Write([]byte("\n"))

	return nil
}
