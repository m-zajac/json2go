package json2go

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
)

func astPrintDecls(decls []ast.Decl) string {
	file := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: decls,
	}

	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), file)

	// remove go file header
	repr := buf.String()
	repr = strings.TrimPrefix(repr, "package main")
	repr = strings.TrimSpace(repr)

	return repr
}
