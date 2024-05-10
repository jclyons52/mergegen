package main

import "go/ast"

func TransformAstToTemplateData(node *ast.File) templateData {
	structs, structNames := collectStructs(node)
	return formatTemplateData(node.Name.Name, structs, structNames)
}

func formatTemplateData(packageName string, structs map[string][]field, structNames []string) templateData {
	data := templateData{
		PackageName: packageName,
	}

	for _, name := range structNames {
		data.Structs = append(data.Structs, structData{
			TypeName: name,
			Fields:   structs[name],
		})
	}

	return data
}

func collectStructs(node *ast.File) (map[string][]field, []string) {
	var structs = make(map[string][]field, 0)
	var structNames []string
	// Collect all struct types
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if st, ok := x.Type.(*ast.StructType); ok {
				typeName := x.Name.Name
				structNames = append(structNames, typeName)
				var fields []field
				for _, f := range st.Fields.List {
					for _, name := range f.Names {
						fieldType, typeElement, isStruct, isPointer := getTypeName(f.Type)
						fields = append(fields, field{
							Name:        name.Name,
							Type:        fieldType,
							TypeElement: typeElement,
							IsStruct:    isStruct,
							IsPointer:   isPointer,
						})
					}
				}
				structs[typeName] = fields
			}
		}
		return true
	})

	return structs, structNames
}

func getTypeName(expr ast.Expr) (string, string, bool, bool) {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name, x.Name, !isBasicType(x.Name), false
	case *ast.SelectorExpr:
		typeName, _, isStruct, isPointer := getTypeName(x.X)
		return typeName + "." + x.Sel.Name, x.Sel.Name, isStruct, isPointer // Preserve struct and pointer flags
	case *ast.StarExpr:
		_, element, isStruct, _ := getTypeName(x.X)
		return "*" + element, element, isStruct, true // Mark as pointer
	case *ast.ArrayType:
		_, element, isStruct, isPointer := getTypeName(x.Elt)
		return "[]" + element, element, isStruct, isPointer // Arrays might be of pointer types
	default:
		return "", "", false, false
	}
}

func isBasicType(name string) bool {
	switch name {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "byte", "rune", "float32", "float64", "complex64", "complex128", "bool", "error":
		return true
	default:
		return false
	}
}
