//go:build js && wasm
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

		rootName := json2go.DefaultRootName
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
	rootName = json2go.DefaultRootName

	if jsVal.Type() != js.TypeObject {
		return nil, rootName
	}

	var useMapsMinAttrs uint = json2go.DefaultMakeMapsWhenMinAttributes
	useMaps := jsVal.Get("useMaps").Truthy()
	if useMaps {
		if v := jsVal.Get("useMapsMinAttrs"); v.Type() == js.TypeNumber {
			useMapsMinAttrs = uint(v.Int())
		} else if v := v.String(); v != "" {
			if w, err := strconv.ParseUint(v, 10, 64); err == nil {
				useMapsMinAttrs = uint(w)
			}
		}
	}

	similarityThreshold := json2go.DefaultSimilarityThreshold
	if v := jsVal.Get("similarityThreshold"); v.Type() == js.TypeNumber {
		similarityThreshold = v.Float()
	} else if v := v.String(); v != "" {
		if w, err := strconv.ParseFloat(v, 64); err == nil {
			similarityThreshold = w
		}
	}

	minSubsetSize := json2go.DefaultMinSubsetSize
	if v := jsVal.Get("minSubsetSize"); v.Type() == js.TypeNumber {
		minSubsetSize = v.Int()
	} else if v := v.String(); v != "" {
		if w, err := strconv.Atoi(v); err == nil {
			minSubsetSize = w
		}
	}

	minSubsetOccurrences := json2go.DefaultMinSubsetOccurrences
	if v := jsVal.Get("minSubsetOccurrences"); v.Type() == js.TypeNumber {
		minSubsetOccurrences = v.Int()
	} else if v := v.String(); v != "" {
		if w, err := strconv.Atoi(v); err == nil {
			minSubsetOccurrences = w
		}
	}

	minAddedFields := json2go.DefaultMinAddedFields
	if v := jsVal.Get("minAddedFields"); v.Type() == js.TypeNumber {
		minAddedFields = v.Int()
	} else if v := v.String(); v != "" {
		if w, err := strconv.Atoi(v); err == nil {
			minAddedFields = w
		}
	}

	opts = append(
		opts,
		json2go.OptExtractAllTypes(
			jsVal.Get("extractAllTypes").Truthy(),
		),
		json2go.OptExtractCommonTypes(
			jsVal.Get("extractCommonTypes").Truthy(),
		),
		json2go.OptExtractHeuristics(
			similarityThreshold,
			minSubsetSize,
			minSubsetOccurrences,
			minAddedFields,
		),
		json2go.OptStringPointersWhenKeyMissing(
			jsVal.Get("stringPointersWhenKeyMissing").Truthy(),
		),
		json2go.OptSkipEmptyKeys(
			jsVal.Get("skipEmptyKeys").Truthy(),
		),
		json2go.OptMakeMaps(
			jsVal.Get("useMaps").Truthy(),
			useMapsMinAttrs,
		),
		json2go.OptTimeAsString(
			jsVal.Get("timeAsStr").Truthy(),
		),
	)

	ro := jsVal.Get("rootName")
	if ro.Type() == js.TypeString {
		if v := ro.String(); v != "" {
			rootName = v
		}
	}

	return opts, rootName
}
