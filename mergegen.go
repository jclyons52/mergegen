package main

import (
	"go/ast"
	"strings"
)

func TransformAstToTemplateData(node *ast.File) templateData {
	structs, structNames, imports := collectStructs(node)
	data := templateData{
		PackageName: node.Name.Name,
		Imports:     deduplicate(imports),
	}

	for _, name := range structNames {
		data.Structs = append(data.Structs, structData{
			TypeName: name,
			Fields:   structs[name],
		})
	}

	return data
}

func collectStructs(node *ast.File) (map[string][]field, []string, []string) {
	var structs = make(map[string][]field, 0)
	var structNames []string
	var imports []string
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
						fieldType, typeElement, isStruct, isPointer, pkg := getTypeName(f.Type)
						fields = append(fields, field{
							Name:        name.Name,
							Type:        fieldType,
							TypeElement: typeElement,
							IsStruct:    isStruct,
							IsPointer:   isPointer,
							IsExternal:  isExternalType(fieldType),
						})
						if isExternalType(fieldType) {
							imports = append(imports, pkg)
						}
					}
				}
				structs[typeName] = fields
			}
		}
		return true
	})

	return structs, structNames, imports
}

func deduplicate(imports []string) []string {
	seen := make(map[string]struct{}, len(imports))
	j := 0
	for _, v := range imports {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		imports[j] = v
		j++
	}
	return imports[:j]
}

func getTypeName(expr ast.Expr) (string, string, bool, bool, string) {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name, x.Name, !isBasicType(x.Name), false, ""
	case *ast.SelectorExpr:
		// Fully qualified name for external types
		pkgIdent, ok := x.X.(*ast.Ident)
		if ok {
			return pkgIdent.Name + "." + x.Sel.Name, x.Sel.Name, false, false, pkgIdent.Name
		}
		typeName, _, isStruct, isPointer, pkg := getTypeName(x.X)
		return typeName + "." + x.Sel.Name, x.Sel.Name, isStruct, isPointer, pkg // Preserve struct and pointer flags
	case *ast.StarExpr:
		_, element, isStruct, _, _ := getTypeName(x.X)
		return "*" + element, element, isStruct, true, "" // Mark as pointer
	case *ast.ArrayType:
		_, element, isStruct, isPointer, _ := getTypeName(x.Elt)
		return "[]" + element, element, isStruct, isPointer, "" // Arrays might be of pointer types
	default:
		return "", "", false, false, ""
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

func isExternalType(typeName string) bool {
	// Simple check if type is considered external or complex enough for mergo
	return strings.Contains(typeName, ".")
}
