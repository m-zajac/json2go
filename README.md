# json2go ![Build](https://github.com/m-zajac/json2go/workflows/Build/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/m-zajac/json2go)](https://goreportcard.com/report/github.com/m-zajac/json2go) [![GoDoc](https://godoc.org/github.com/m-zajac/json2go?status.svg)](http://godoc.org/github.com/m-zajac/json2go) [![Coverage](https://img.shields.io/badge/coverage-gocover.io-blue)](https://gocover.io/github.com/m-zajac/json2go) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

Package json2go provides utilities for creating go type representation from json inputs.

Json2go can be used in various ways:

- CLI tool
- Web page with conversion tool: [https://m-zajac.github.io/json2go](https://m-zajac.github.io/json2go)
- Go package
- VSCode extension: [vsc-json2go](https://marketplace.visualstudio.com/items?itemName=m-zajac.vsc-json2go)

## Why another conversion tool?

There are few tools for converting json to go types available already. But all of which I tried worked correctly with only basic documents. 

The goal of this project is to create types, that are guaranteed to properly unmarshal input data. There are multiple test cases that check marshaling/unmarshaling both ways to ensure generated type is accurate.

Here's the micro acid test if you want to check this or other conversion tools:

```json
[{"date":"2020-10-03T15:04:05Z","text":"txt1","doc":{"x":"x"}},{"date":"2020-10-03T15:05:02Z","_doc":false,"arr":[[1,null,1.23]]},{"bool":true,"doc":{"y":123}}]
```

And correct output with some comments:

```go
type Document []struct {
	Arr  [][]*float64 `json:"arr,omitempty"` // Should be doubly nested array; should be a pointer type because there's null in values.
	Bool *bool        `json:"bool,omitempty"` // Shouldn't be `bool` because when key is missing you'll get false information.
	Date *time.Time   `json:"date,omitempty"` // Could be also `string` or `*string`.
	Doc  *Doc         `json:"doc,omitempty"`  // Should be pointer, because key is not present in all documents in array.
	Doc2 *bool        `json:"_doc,omitempty"` // Attribute for "_doc" key (other that for "doc"!). Type - the same as `Bool` attribute.
	Text *string      `json:"text,omitempty"` // Could be also `string`.
}

type Doc struct {
	X *string `json:"x,omitempty"` // Should be pointer, because key is not present in all documents.
	Y *int    `json:"y,omitempty"` // Should be pointer, because key is not present in all documents.
}
```

## CLI Installation

    go install github.com/m-zajac/json2go/cmd/json2go@latest

## Usage

Json2go can be used as cli tool or as package.

CLI tools can be used directly to create go type from stdin data (see examples).

Package provides Parser, which can consume multiple jsons and outputs go type fitting all inputs (see examples and [documentation](https://godoc.org/github.com/m-zajac/json2go)). Example usage: read documents from document-oriented database and feed them too parser for go struct.

### CLI usage examples

    echo '{"x":1,"y":2}' | json2go

---

    curl -s https://api.punkapi.com/v2/beers?page=1&per_page=5 | json2go

Check this one :)

---

    cat data.json | json2go

### Package usage examples

```go
inputs := []string{
	`{"x": 123, "y": "test", "z": false}`,
	`{"a": 123, "x": 12.3, "y": true}`,
}

parser := json2go.NewJSONParser("Document")
for _, in := range inputs {
	parser.FeedBytes([]byte(in))
}

res := parser.String()
fmt.Println(res)
```

## Example outputs

```json
{
    "line": {
        "point1": {
            "x": 12.1,
            "y": 2
        },
        "point2": {
            "x": 12.1,
            "y": 2
        }
    }
}
```
```go
type Document struct {
	Line Line `json:"line"`
}

type Point struct {
	X float64 `json:"x"`
	Y int     `json:"y"`
}

type Line struct {
	Point1 Point `json:"point1"`
	Point2 Point `json:"point2"`
}
```

---

```json
[
    {
        "name": "water",
        "type": "liquid",
        "boiling_point": {
            "units": "C",
            "value": 100
        }
    },
    {
        "name": "oxygen",
        "type": "gas",
        "density": {
            "units": "g/L",
            "value": 1.429
        }
    },
    {
        "name": "carbon monoxide",
        "type": "gas",
        "dangerous": true,
        "boiling_point": {
            "units": "C",
            "value": -191.5
        },
        "density": {
            "units": "kg/m3",
            "value": 789
        }
    }
]
```
```go
type Document []struct {
	BoilingPoint *UnitsValue `json:"boiling_point,omitempty"`
	Dangerous    *bool       `json:"dangerous,omitempty"`
	Density      *UnitsValue `json:"density,omitempty"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
}

type UnitsValue struct {
	Units string  `json:"units"`
	Value float64 `json:"value"`
}
```

## Contribution

I'd love your input! Especially reporting bugs and proposing new features.
