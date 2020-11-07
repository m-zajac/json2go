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

func astMakeDecls(rootNodes []*node, opts options) []ast.Decl {
	var decls []ast.Decl

	for _, node := range rootNodes {
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(node.name),
					Type: astTypeFromNode(node, opts),
				},
			},
		})
	}

	return decls
}

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

func astTypeFromNode(n *node, opts options) ast.Expr {
	var resultType ast.Expr
	notRequiredAsPointer := true
	allowPointer := true

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
		resultType = astStructTypeFromNode(n, opts)
	case nodeExtractedType:
		resultType = astTypeFromExtractedNode(n)
	case nodeInterfaceType, nodeInitType:
		resultType = newEmptyInterfaceExpr()
		allowPointer = false
	case nodeMapType:
		resultType = astTypeFromMapNode(n, opts)
		allowPointer = false
	default:
		panic(fmt.Sprintf("unknown type: %v", n.t))
	}

	if astTypeShouldBeAPointer(n, notRequiredAsPointer, allowPointer) {
		resultType = &ast.StarExpr{
			X: resultType,
		}
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

func astTypeFromMapNode(n *node, opts options) ast.Expr {
	var ve ast.Expr
	if len(n.children) == 0 {
		ve = newEmptyInterfaceExpr()
	} else {
		ve = astTypeFromNode(n.children[0], opts)
	}
	return &ast.MapType{
		Key:   ast.NewIdent("string"),
		Value: ve,
	}
}

func astTypeFromExtractedNode(n *node) ast.Expr {
	extName := n.externalTypeID
	if extName == "" {
		extName = n.name
	}
	return ast.NewIdent(extName)
}

func astStructTypeFromNode(n *node, opts options) *ast.StructType {
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
		typeDesc.Fields.List = append(typeDesc.Fields.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(child.name)},
			Type:  astTypeFromNode(child.node, opts),
			Tag:   astJSONTag(child.node.key, !child.node.required),
		})
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

func astTypeShouldBeAPointer(n *node, notRequiredAsPointer bool, allowPointer bool) bool {
	if !allowPointer {
		return false
	}

	if !n.root && n.arrayLevel == 0 {
		if n.nullable || (!n.required && notRequiredAsPointer) {
			return true
		}
	} else if n.arrayLevel > 0 {
		if n.arrayWithNulls {
			return true
		}
	}

	return false
}

func newEmptyInterfaceExpr() ast.Expr {
	return &ast.InterfaceType{
		Methods: &ast.FieldList{
			Opening: token.Pos(1),
			Closing: token.Pos(2),
		},
	}
}
