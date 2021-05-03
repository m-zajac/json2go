// +build js,wasm

package main

import (
	"encoding/json"
	"strconv"
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
			parserOpts, rootName = parseOpts(args[1])
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

func parseOpts(jsVal js.Value) (opts []json2go.JSONParserOpt, rootName string) {
	rootName = "document"

	if jsVal.Type() != js.TypeObject {
		return nil, rootName
	}

	if v := jsVal.Get("extractCommonTypes"); !v.IsUndefined() {
		opts = append(opts, json2go.OptExtractCommonTypes(v.Truthy()))
	}
	if v := jsVal.Get("stringPointersWhenKeyMissing"); !v.IsUndefined() {
		opts = append(opts, json2go.OptStringPointersWhenKeyMissing(v.Truthy()))
	}
	if v := jsVal.Get("skipEmptyKeys"); !v.IsUndefined() {
		opts = append(opts, json2go.OptSkipEmptyKeys(v.Truthy()))
	}
	if v := jsVal.Get("useMaps"); !v.IsUndefined() {
		var useMapsMinAttrs uint = 5
		if v := jsVal.Get("useMapsMinAttrs").String(); v != "" {
			if w, err := strconv.ParseUint(v, 10, 64); err == nil {
				useMapsMinAttrs = uint(w)
			}
		}
		opts = append(opts, json2go.OptMakeMaps(v.Truthy(), useMapsMinAttrs))
	}
	if v := jsVal.Get("timeAsStr"); !v.IsUndefined() {
		opts = append(opts, json2go.OptTimeAsString(v.Truthy()))
	}
	if v := jsVal.Get("findRecurrence"); !v.IsUndefined() {
		opts = append(opts, json2go.OptFindRecurrence(v.Truthy()))
	}

	ro := jsVal.Get("rootName")
	if ro.Type() == js.TypeString {
		if v := ro.String(); v != "" {
			rootName = v
		}
	}

	return opts, rootName
}
