package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	htmlFS := http.FileServer(http.Dir("."))
	wasmFS := http.FileServer(http.Dir("../../build/web"))

	log.Print("Serving on http://localhost:8080")
	_ = http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Cache-Control", "no-cache")
		switch {
		case strings.HasSuffix(req.URL.Path, ".wasm"):
			resp.Header().Set("content-type", "application/wasm")
			wasmFS.ServeHTTP(resp, req)
		case strings.HasSuffix(req.URL.Path, "wasm_exec.js"):
			wasmFS.ServeHTTP(resp, req)
		default:
			htmlFS.ServeHTTP(resp, req)
		}
	}))
}
