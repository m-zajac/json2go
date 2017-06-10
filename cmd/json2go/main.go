package main

import (
	"fmt"
	"os"

	"github.com/m-zajac/json2go"
)

func main() {
	if err := json2go.Decode(os.Stdin, os.Stdout); err != nil {
		fmt.Printf("%v\n", err)
	}
}
