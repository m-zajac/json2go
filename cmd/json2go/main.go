package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/m-zajac/json2go"
)

func main() {
	extractCommonNodes := flag.Bool("c", json2go.DefaultExtractCommonTypes, "Extract common nodes as top level struct definitions")
	similarityThreshold := flag.Float64("ct", json2go.DefaultSimilarityThreshold, "Minimum similarity score (0.0 to 1.0) to merge types")
	minSubsetSize := flag.Int("cs", json2go.DefaultMinSubsetSize, "Minimum number of fields in a shared subset to extract it")
	minSubsetOccurrences := flag.Int("co", json2go.DefaultMinSubsetOccurrences, "Minimum frequency of a subset to be extracted")
	minAddedFields := flag.Int("cf", json2go.DefaultMinAddedFields, "Minimum unique fields required for a new type")
	stringPointers := flag.Bool("sp", json2go.DefaultStringPointersWhenKeyMissing, "Allow string pointers when string key is missing in one of documents")
	skipEmptyKeys := flag.Bool("k", json2go.DefaultSkipEmptyKeys, "Ignore keys that were only nulls")
	useMaps := flag.Bool("m", json2go.DefaultMakeMaps, "Try to use maps instead of structs where possible")
	useMapsMinAttrs := flag.Int("mk", json2go.DefaultMakeMapsWhenMinAttributes, "Minimum number of attributes in object to try converting it to a map.")
	timeAsStr := flag.Bool("st", json2go.DefaultTimeAsStr, "Don't use time.Time type, just strings")
	rootTypeName := flag.String("n", json2go.DefaultRootName, "Type name")

	flag.Parse()

	var data interface{}

	jd := json.NewDecoder(os.Stdin)
	if err := jd.Decode(&data); err != nil {
		log.Fatalf("json decoding error: %v", err)
	}

	parser := json2go.NewJSONParser(
		*rootTypeName,
		json2go.OptExtractCommonTypes(*extractCommonNodes),
		json2go.OptExtractHeuristics(*similarityThreshold, *minSubsetSize, *minSubsetOccurrences, *minAddedFields),
		json2go.OptStringPointersWhenKeyMissing(*stringPointers),
		json2go.OptSkipEmptyKeys(*skipEmptyKeys),
		json2go.OptMakeMaps(*useMaps, uint(*useMapsMinAttrs)),
		json2go.OptTimeAsString(*timeAsStr),
	)

	parser.FeedValue(data)

	repr := parser.String()

	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(repr)
	os.Stdout.WriteString("\n\n")
}
