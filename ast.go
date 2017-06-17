package json2go

import (
	"bytes"
	"go/ast"
	"go/token"
	"sort"
)

func astMakeDecls(rootNode *node) []ast.Decl {
	name := attrName(rootNode.name)
	if name == "" {
		return nil
	}
	rootDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(name),
				Type: astTypeFromNode(rootNode),
			},
		},
	}

	return []ast.Decl{rootDecl}
}

func astTypeFromNode(n *node) ast.Expr {
	switch n.t.id {
	case nodeTypeBool:
		return ast.NewIdent("bool")
	case nodeTypeInt:
		return ast.NewIdent("int")
	case nodeTypeFloat:
		return ast.NewIdent("float64")
	case nodeTypeString:
		return ast.NewIdent("string")

	case nodeTypeArrayUnknown:
		return &ast.ArrayType{
			Elt: &ast.InterfaceType{
				Methods: &ast.FieldList{
					Opening: 1,
					Closing: 2,
				},
			},
		}
	case nodeTypeArrayBool:
		return &ast.ArrayType{
			Elt: ast.NewIdent("bool"),
		}
	case nodeTypeArrayInt:
		return &ast.ArrayType{
			Elt: ast.NewIdent("int"),
		}
	case nodeTypeArrayFloat:
		return &ast.ArrayType{
			Elt: ast.NewIdent("float64"),
		}
	case nodeTypeArrayString:
		return &ast.ArrayType{
			Elt: ast.NewIdent("string"),
		}
	case nodeTypeArrayInterface:
		return &ast.ArrayType{
			Elt: &ast.InterfaceType{
				Methods: &ast.FieldList{
					Opening: 1,
					Closing: 2,
				},
			},
		}

	case nodeTypeArrayObject:
		return &ast.ArrayType{
			Elt: astStructTypeFromNode(n),
		}

	case nodeTypeObject:
		return astStructTypeFromNode(n)
	}

	return &ast.InterfaceType{
		Methods: &ast.FieldList{
			Opening: token.Pos(1),
			Closing: token.Pos(2),
		},
	}
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
	for childName, child := range n.children {
		sortedChildren = append(sortedChildren, nodeWithName{
			key:  childName,
			node: child,
		})
	}
	sort.Slice(sortedChildren, func(i, j int) bool {
		return sortedChildren[i].key < sortedChildren[j].key
	})

	for _, child := range sortedChildren {
		childName := attrName(child.node.name)
		if childName == "" {
			continue
		}

		typeDesc.Fields.List = append(typeDesc.Fields.List, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(childName)},
			Type:  astTypeFromNode(child.node),
			Tag:   astJSONTag(child.key, false),
		})
	}

	return typeDesc
}

func astJSONTag(name string, omitempty bool) *ast.BasicLit {
	var buf bytes.Buffer
	buf.WriteString("`json:\"")
	buf.WriteString(name)
	if omitempty {
		buf.WriteString(",omitempty")
	}
	buf.WriteString("\"`")

	return &ast.BasicLit{
		Value: buf.String(),
	}
}
