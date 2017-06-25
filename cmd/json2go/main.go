package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/m-zajac/json2go"
)

func main() {
	extractCommonNodes := flag.Bool("c", true, "Extract common nodes as top level struct definitions")
	rootTypeName := flag.String("n", "Document", "Type name")

	flag.Parse()

	var data interface{}

	jd := json.NewDecoder(os.Stdin)
	if err := jd.Decode(&data); err != nil {
		log.Fatalf("json decoding error: %v", err)
	}

	parser := json2go.NewJSONParser(*rootTypeName)
	parser.ExtractCommonStructs = *extractCommonNodes

	parser.FeedValue(data)

	repr := parser.String()

	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(repr)
	os.Stdout.WriteString("\n\n")
}
