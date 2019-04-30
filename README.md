# json2go [![Build Status](https://travis-ci.org/m-zajac/json2go.svg?branch=master)](https://travis-ci.org/m-zajac/json2go) [![Go Report Card](https://goreportcard.com/badge/github.com/m-zajac/json2go)](https://goreportcard.com/report/github.com/m-zajac/json2go) [![GoDoc](https://godoc.org/github.com/m-zajac/json2go?status.svg)](http://godoc.org/github.com/m-zajac/json2go)

Package json2go provides utilities for creating well formated go type representation from json inputs.

## Demo page

Simple web page with webassembly json2go parser

[https://m-zajac.github.io/json2go](https://m-zajac.github.io/json2go)

## Installation

    go get github.com/m-zajac/json2go/...

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
	Line struct {
		Point1 Point `json:"point1"`
		Point2 Point `json:"point2"`
	} `json:"line"`
}
type Point struct {
	X float64 `json:"x"`
	Y int     `json:"y"`
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


## TODO

- extract common types (working, but every nested struct had to have same fields)
- try decoding to map if resulting struct has many attributes with same type
- convert json schema to go type
