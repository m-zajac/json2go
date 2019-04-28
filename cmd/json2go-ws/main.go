// +build js,wasm
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/m-zajac/json2go"
)

func main() {
	in := js.Global().Get("document").Call("getElementById", "input")
	out := js.Global().Get("document").Call("getElementById", "out")
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := in.Get("value").String()

		var data interface{}
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			out.Set(
				"innerHTML",
				fmt.Sprintf("json decoding error: %v", err),
			)
			return nil
		}

		parser := json2go.NewJSONParser("root")
		parser.ExtractCommonStructs = true
		parser.FeedValue(data)

		out.Set("innerHTML", parser.String())

		return nil
	})

	js.Global().Get("document").Call("getElementById", "run").Call("addEventListener", "click", cb)

	select {}
}
