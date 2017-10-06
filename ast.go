package json2go

import (
	"bytes"
	"go/ast"
	"go/token"
	"sort"
)

func astMakeDecls(rootNodes []*node) []ast.Decl {
	var decls []ast.Decl

	for _, rootNode := range rootNodes {
		name := attrName(rootNode.key)
		if name == "" {
			continue
		}
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(name),
					Type: astTypeFromNode(rootNode),
				},
			},
		})
	}

	return decls
}

func astTypeFromNode(n *node) ast.Expr {
	var resultType ast.Expr
	var pointable bool

	switch n.t.(type) {
	case nodeBoolType:
		resultType = ast.NewIdent("bool")
		pointable = true
	case nodeIntType:
		resultType = ast.NewIdent("int")
		pointable = true
	case nodeFloatType:
		resultType = ast.NewIdent("float64")
		pointable = true
	case nodeStringType:
		resultType = ast.NewIdent("string")
		pointable = false
	case nodeObjectType:
		resultType = astStructTypeFromNode(n)
		pointable = true
	case nodeExternalType:
		extName := n.externalTypeID
		if extName == "" {
			extName = attrName(n.key)
		}
		resultType = ast.NewIdent(extName)
		pointable = true
	default:
		resultType = &ast.InterfaceType{
			Methods: &ast.FieldList{
				Opening: token.Pos(1),
				Closing: token.Pos(2),
			},
		}
	}

	if pointable && !n.required && !n.root && n.arrayLevel == 0 {
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

func astStructTypeFromNode(n *node) *ast.StructType {
	typeDesc := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{},
		},
	}

	// sort children by name
	type nodeWithName struct {
		key  string
		node *node
	}
	var sortedChildren []nodeWithName
	for _, child := range n.children {
		sortedChildren = append(sortedChildren, nodeWithName{
			key:  child.key,
			node: child,
		})
	}
	sort.Slice(sortedChildren, func(i, j int) bool {
		return sortedChildren[i].key < sortedChildren[j].key
	})

	for _, child := range sortedChildren {
		childName := attrName(child.key)
		if childName == "" {
			continue
		}

		typeDesc.Fields.List = append(typeDesc.Fields.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(childName)},
			Type:  astTypeFromNode(child.node),
			Tag:   astJSONTag(child.key, !child.node.required),
		})
	}

	return typeDesc
}

func astJSONTag(key string, omitempty bool) *ast.BasicLit {
	var buf bytes.Buffer
	buf.WriteString("`json:\"")
	buf.WriteString(key)
	if omitempty {
		buf.WriteString(",omitempty")
	}
	buf.WriteString("\"`")

	return &ast.BasicLit{
		Value: buf.String(),
	}
}
