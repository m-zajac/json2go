# json2go [![Build Status](https://travis-ci.org/m-zajac/json2go.svg?branch=master)](https://travis-ci.org/m-zajac/json2go) [![Go Report Card](https://goreportcard.com/badge/github.com/m-zajac/json2go)](https://goreportcard.com/report/github.com/m-zajac/json2go) [![GoDoc](https://godoc.org/github.com/m-zajac/json2go?status.svg)](http://godoc.org/github.com/m-zajac/json2go)

Package json2go provides utilities for creating well formated go type representation from json inputs.

## Installation

    go get github.com/m-zajac/json2go/...

## Usage

Json2go can be used as cli tool or as package.

CLI tools can be used directly to create go type from stdin data (see examples).

Package provides Parser, which can consume multiple jsons and outputs go type fitting all inputs (see examples and [documentation](https://godoc.org/github.com/m-zajac/json2go)). Example usage: read documents from document-oriented database and feed them too parser for go struct.

### CLI examples:

    echo '123' | json2go

    >>>

    type Object int

---

    echo '{"x": 123, "y": "test", "z": false}' | json2go

    >>>

    type Object struct {
    	X	int      `json:"x"`
    	Y	string   `json:"y"`
    	Z	bool     `json:"z"`
    }

---

    echo '[{"x": 123, "y": "test", "z": false}, {"a": 123, "x": 12.3, "y": true}]' | json2go

    >>>

    type Object struct {
    	A	int          `json:"a"`
    	X	float64      `json:"x"`
    	Y	interface{}  `json:"y"`
    	Z	bool         `json:"z"`
    }

---

    curl -s https://www.reddit.com/r/golang.json | json2go

    >>>

    Check out yourself :)

### Package examples:

```go
inputs = []string{
	`{"x": 123, "y": "test", "z": false}`,
	`{"a": 123, "x": 12.3, "y": true}`,
}

parser := json2go.NewParser()
for _, in := range inputs {
	parser.FeedBytes([]byte(in))
}

err, res := parser.Result()
if err != nil {
	log.Errorf("parser error: %v", err)
}
fmt.Println(res)
```


## TODO

- extract common structures
