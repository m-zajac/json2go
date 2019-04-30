// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/m-zajac/json2go"
)

func main() {
	js.Global().Set("json2go", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return ""
		}
		input := args[0].String()

		parser := json2go.NewJSONParser("document")
		parser.ExtractCommonTypes = true

		if len(args) > 1 {
			opts := args[1]
			if opts.Type() == js.TypeObject {
				parser.ExtractCommonTypes = opts.Get("extractCommonTypes").Truthy()
			}
		}

		var data interface{}
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			return ""
		}

		parser.FeedValue(data)

		return parser.String()
	}))

	select {}
}
