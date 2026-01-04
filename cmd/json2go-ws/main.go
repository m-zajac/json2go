//go:build js && wasm

package main

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/m-zajac/json2go"
	"github.com/tidwall/gjson"
)

func main() {
	json2goFunc := js.FuncOf(json2goHandler)
	js.Global().Set("json2go", json2goFunc)

	select {}
}

// json2goHandler handles calls from JavaScript to convert JSON to Go structs.
// It expects 1-2 arguments:
//   - args[0]: JSON string to parse (required)
//   - args[1]: Options object (optional)
//
// Returns: Generated Go struct code as a string, or empty string on error
func json2goHandler(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return ""
	}

	if args[0].Type() != js.TypeString {
		return ""
	}

	input := args[0].String()
	if input == "" {
		return ""
	}

	inputBytes := []byte(input)
	if !gjson.ValidBytes(inputBytes) {
		return ""
	}

	rootName := json2go.DefaultRootName
	var parserOpts []json2go.JSONParserOpt
	var err error

	if len(args) > 1 {
		parserOpts, rootName, err = parseOpts(args[1])
		if err != nil {
			return ""
		}
	}

	data := gjson.ParseBytes(inputBytes).Value()

	parser := json2go.NewJSONParser(rootName, parserOpts...)
	parser.FeedValue(data)

	return parser.String()
}

func parseOpts(jsVal js.Value) (opts []json2go.JSONParserOpt, rootName string, err error) {
	rootName = json2go.DefaultRootName

	// If not an object, return defaults (not an error - options are optional)
	if jsVal.Type() != js.TypeObject {
		return nil, rootName, nil
	}

	// Parse useMapsMinAttrs (minimum attributes for map detection)
	useMapsMinAttrs := uint(json2go.DefaultMakeMapsWhenMinAttributes)
	useMaps := jsVal.Get("useMaps").Truthy()
	if useMaps {
		if v, err := parseUintOption(jsVal, "useMapsMinAttrs", 1); err != nil {
			return nil, "", fmt.Errorf("useMapsMinAttrs: %w", err)
		} else if v != 0 {
			useMapsMinAttrs = v
		}
	}

	// Parse similarityThreshold (0.0 to 1.0)
	similarityThreshold := json2go.DefaultSimilarityThreshold
	if v, err := parseFloatOption(jsVal, "similarityThreshold"); err != nil {
		return nil, "", fmt.Errorf("similarityThreshold: %w", err)
	} else if v != 0 {
		if v < 0.0 || v > 1.0 {
			return nil, "", fmt.Errorf("similarityThreshold must be between 0.0 and 1.0, got %f", v)
		}
		similarityThreshold = v
	}

	// Parse minSubsetSize (minimum 1)
	minSubsetSize := json2go.DefaultMinSubsetSize
	if v, err := parseIntOption(jsVal, "minSubsetSize", 1); err != nil {
		return nil, "", fmt.Errorf("minSubsetSize: %w", err)
	} else if v != 0 {
		minSubsetSize = v
	}

	// Parse minSubsetOccurrences (minimum 1)
	minSubsetOccurrences := json2go.DefaultMinSubsetOccurrences
	if v, err := parseIntOption(jsVal, "minSubsetOccurrences", 1); err != nil {
		return nil, "", fmt.Errorf("minSubsetOccurrences: %w", err)
	} else if v != 0 {
		minSubsetOccurrences = v
	}

	// Parse minAddedFields (minimum 1)
	minAddedFields := json2go.DefaultMinAddedFields
	if v, err := parseIntOption(jsVal, "minAddedFields", 1); err != nil {
		return nil, "", fmt.Errorf("minAddedFields: %w", err)
	} else if v != 0 {
		minAddedFields = v
	}

	// Build parser options
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

	// Parse root name
	if ro := jsVal.Get("rootName"); ro.Type() == js.TypeString {
		if v := ro.String(); v != "" {
			rootName = v
		}
	}

	return opts, rootName, nil
}

func parseIntOption(jsVal js.Value, key string, minimum int) (int, error) {
	v := jsVal.Get(key)

	// Not set - return 0 to signal "use default"
	if v.IsUndefined() || v.IsNull() {
		return 0, nil
	}

	var result int

	if v.Type() == js.TypeNumber {
		result = v.Int()
	} else if v.Type() == js.TypeString {
		s := v.String()
		if s == "" {
			return 0, nil
		}
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("invalid integer value %q", s)
		}
		result = parsed
	} else {
		return 0, fmt.Errorf("must be a number or string, got %s", v.Type().String())
	}

	if result < minimum {
		return 0, fmt.Errorf("must be at least %d, got %d", minimum, result)
	}

	return result, nil
}

func parseUintOption(jsVal js.Value, key string, minimum uint) (uint, error) {
	v := jsVal.Get(key)

	// Not set - return 0 to signal "use default"
	if v.IsUndefined() || v.IsNull() {
		return 0, nil
	}

	var result uint

	if v.Type() == js.TypeNumber {
		intVal := v.Int()
		if intVal < 0 {
			return 0, fmt.Errorf("must be non-negative, got %d", intVal)
		}
		result = uint(intVal)
	} else if v.Type() == js.TypeString {
		s := v.String()
		if s == "" {
			return 0, nil
		}
		parsed, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid unsigned integer value %q", s)
		}
		result = uint(parsed)
	} else {
		return 0, fmt.Errorf("must be a number or string, got %s", v.Type().String())
	}

	if result < minimum {
		return 0, fmt.Errorf("must be at least %d, got %d", minimum, result)
	}

	return result, nil
}

func parseFloatOption(jsVal js.Value, key string) (float64, error) {
	v := jsVal.Get(key)

	// Not set - return 0 to signal "use default"
	if v.IsUndefined() || v.IsNull() {
		return 0.0, nil
	}

	if v.Type() == js.TypeNumber {
		return v.Float(), nil
	} else if v.Type() == js.TypeString {
		s := v.String()
		if s == "" {
			return 0.0, nil
		}
		parsed, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0.0, fmt.Errorf("invalid float value %q", s)
		}
		return parsed, nil
	}

	return 0.0, fmt.Errorf("must be a number or string, got %s", v.Type().String())
}
