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

		rootName := "document"
		var parserOpts []json2go.JSONParserOpt
		if len(args) > 1 {
			opts := args[1]
			if opts.Type() == js.TypeObject {
				parserOpts = append(
					parserOpts,
					json2go.OptExtractCommonTypes(
						opts.Get("extractCommonTypes").Truthy(),
					),
					json2go.OptStringPointersWhenKeyMissing(
						opts.Get("stringPointersWhenKeyMissing").Truthy(),
					),
					json2go.OptSkipEmptyKeys(
						opts.Get("skipEmptyKeys").Truthy(),
					),
				)

				ro := opts.Get("rootName")
				if ro.Type() == js.TypeString {
					if v := ro.String(); v != "" {
						rootName = v
					}
				}
			}
		}

		parser := json2go.NewJSONParser(rootName, parserOpts...)

		var data interface{}
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			return ""
		}

		parser.FeedValue(data)

		return parser.String()
	}))

	select {}
}
