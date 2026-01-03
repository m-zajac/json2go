package json2go

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"sort"
	"strings"
)

// Printer settings - copy from gofmt.
// See: https://github.com/golang/go/blob/go1.15.5/src/cmd/gofmt/gofmt.go
const (
	tabWidth    = 8
	printerMode = printer.UseSpaces | printer.TabIndent | printerNormalizeNumbers

	// printerNormalizeNumbers means to canonicalize number literal prefixes
	// and exponents while printing. See https://golang.org/doc/go1.13#gofmt.
	//
	// This value is defined in go/printer specifically for go/format and cmd/gofmt.
	printerNormalizeNumbers = 1 << 30
)

func astMakeDecls(rootNodes []*node, opts options) []ast.Decl {
	var decls []ast.Decl

	for _, node := range rootNodes {
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(node.name),
					Type: astTypeFromNode(node, opts, node.name),
				},
			},
		})
	}

	return decls
}

func astPrintDecls(decls []ast.Decl) string {
	if len(decls) == 0 {
		return ""
	}

	// Print each declaration separately and join with blank lines
	var parts []string
	fset := token.NewFileSet()
	// Use go/printer with settings compatible with gofmt.
	prn := printer.Config{Mode: printerMode, Tabwidth: tabWidth}

	for _, decl := range decls {
		var buf bytes.Buffer
		prn.Fprint(&buf, fset, decl)
		parts = append(parts, strings.TrimSpace(buf.String()))
	}

	return strings.Join(parts, "\n\n")
}

func astTypeFromNode(n *node, opts options, rootNodeName string) ast.Expr {
	var resultType ast.Expr
	notRequiredAsPointer := true

	switch n.t.(type) {
	case nodeBoolType:
		resultType = ast.NewIdent("bool")
	case nodeIntType:
		resultType = ast.NewIdent("int")
	case nodeFloatType:
		resultType = ast.NewIdent("float64")
	case nodeStringType:
		resultType = ast.NewIdent("string")
		notRequiredAsPointer = opts.stringPointersWhenKeyMissing
	case nodeTimeType:
		resultType = astTypeFromTimeNode(n, opts)
		if opts.timeAsStr {
			notRequiredAsPointer = opts.stringPointersWhenKeyMissing
		}
	case nodeObjectType:
		resultType = astStructTypeFromNode(n, opts, rootNodeName)
	case nodeExtractedType:
		resultType = astTypeFromExtractedNode(n)
	case nodeInterfaceType, nodeInitType:
		resultType = newEmptyInterfaceExpr()
	case nodeMapType:
		resultType = astTypeFromMapNode(n, opts, rootNodeName)
	default:
		panic(fmt.Sprintf("unknown type: %v", n.t))
	}

	if astTypeShouldBeAPointer(resultType, n, rootNodeName, notRequiredAsPointer) {
		resultType = &ast.StarExpr{X: resultType}
	}

	for i := n.arrayLevel; i > 0; i-- {
		resultType = &ast.ArrayType{
			Elt: resultType,
		}
	}

	return resultType
}

func astTypeFromTimeNode(n *node, opts options) ast.Expr {
	var resultType ast.Expr

	if opts.timeAsStr {
		resultType = ast.NewIdent("string")
	} else if n.root {
		// We have to use type alias here to preserve "UnmarshalJSON" method from time type.
		resultType = ast.NewIdent("= time.Time")
	} else {
		resultType = ast.NewIdent("time.Time")
	}

	return resultType
}

func astTypeFromMapNode(n *node, opts options, rootNodeName string) ast.Expr {
	var ve ast.Expr
	if len(n.children) == 0 {
		ve = newEmptyInterfaceExpr()
	} else {
		ve = astTypeFromNode(n.children[0], opts, rootNodeName)
	}
	return &ast.MapType{
		Key:   ast.NewIdent("string"),
		Value: ve,
	}
}

func astTypeFromExtractedNode(n *node) ast.Expr {
	extName := n.extractedTypeName
	if extName == "" {
		extName = n.name
	}
	return ast.NewIdent(extName)
}

func astStructTypeFromNode(n *node, opts options, rootNodeName string) *ast.StructType {
	typeDesc := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

	// sort children by name
	type nodeWithName struct {
		name string
		node *node
	}
	var sortedChildren []nodeWithName
	for _, child := range n.children {
		sortedChildren = append(sortedChildren, nodeWithName{
			name: child.name,
			node: child,
		})
	}
	sort.Slice(sortedChildren, func(i, j int) bool {
		return sortedChildren[i].name < sortedChildren[j].name
	})

	for _, child := range sortedChildren {
		field := &ast.Field{
			Type: astTypeFromNode(child.node, opts, rootNodeName),
		}
		if !child.node.embedded {
			field.Names = []*ast.Ident{ast.NewIdent(child.name)}
			field.Tag = astJSONTag(child.node.key, !child.node.required)
		}

		typeDesc.Fields.List = append(typeDesc.Fields.List, field)
	}

	return typeDesc
}

func astJSONTag(key string, omitempty bool) *ast.BasicLit {
	tag := fmt.Sprintf("%#v", key)
	tag = strings.Trim(tag, `"`)
	if omitempty {
		tag = fmt.Sprintf("`json:\"%s,omitempty\"`", tag)
	} else {
		tag = fmt.Sprintf("`json:\"%s\"`", tag)
	}

	return &ast.BasicLit{
		Value: tag,
	}
}

func astTypeShouldBeAPointer(resultType ast.Expr, n *node, rootNodeName string, notRequiredAsPointer bool) bool {
	switch resultType.(type) {
	case *ast.StarExpr:
		// Already a pointer.
		return false
	case *ast.InterfaceType:
		// An empty interface will never be a pointer.
		return false
	case *ast.MapType:
		// Maps will never be a pointer.
		return false
	default:
		if n.arrayLevel == 0 {
			if _, ok := n.t.(nodeExtractedType); ok && (n.extractedTypeName == rootNodeName || (n.extractedTypeName == "" && n.name == rootNodeName)) {
				return true // Recursive reference
			}
			return !n.root && (n.nullable || (!n.required && notRequiredAsPointer))
		}

		return n.arrayWithNulls
	}

}

func newEmptyInterfaceExpr() ast.Expr {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{
			Opening: token.Pos(1),
			Closing: token.Pos(2),
		},
	}
}
