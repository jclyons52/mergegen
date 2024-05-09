package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	srcFile  = flag.String("src", "", "Source file containing the struct definitions.")
	typeName = flag.String("type", "", "Type of the struct to generate the merger for.")
	output   = flag.String("output", "merge_functions.go", "Output file for the generated merge functions.")
)

const funcTemplate = `package {{ .PackageName }}

// generated code, do not modify

{{ range .Structs }}
// Merge{{ .TypeName }} merges two {{ .TypeName }} structs.
func Merge{{ .TypeName }}(dst, src *{{ .TypeName }}) {
{{- range .Fields }}
    {{- if .IsStruct }}
    if dst.{{ .Name }} == nil {
        dst.{{ .Name }} = new({{ .TypeElement }})
    }
    Merge{{ .TypeElement }}(dst.{{ .Name }}, src.{{ .Name }})
    {{- else }}
    if src.{{ .Name }} != {{ defaultZeroValue .Type }} {
        dst.{{ .Name }} = src.{{ .Name }}
    }
    {{- end }}
{{- end }}
}
{{ end }}
`

type field struct {
	Name        string
	Type        string
	TypeElement string
	IsStruct    bool
}

func main() {
	flag.Parse()

	if *srcFile == "" || *typeName == "" || *output == "" {
		fmt.Println("Source file, type name, and output file must be provided.")
		return
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *srcFile, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Failed to parse source file:", err)
		return
	}

	packageName := node.Name.Name // Get the package name from the AST

	structs := make(map[string][]field)
	structNames := make([]string, 0)

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
						fieldType, typeElement, isStruct := getTypeName(f.Type)
						fields = append(fields, field{
							Name:        name.Name,
							Type:        fieldType,
							TypeElement: typeElement,
							IsStruct:    isStruct,
						})
					}
				}
				structs[typeName] = fields
			}
		}
		return true
	})

	funcMap := template.FuncMap{
		"defaultZeroValue": func(typeStr string) string {
			if strings.Contains(typeStr, "int") || strings.Contains(typeStr, "float") {
				return "0"
			}
			if strings.Contains(typeStr, "string") {
				return `""`
			}
			if strings.Contains(typeStr, "bool") {
				return "false"
			}
			return "nil"
		},
	}

	t, err := template.New("func").Funcs(funcMap).Parse(funcTemplate)
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return
	}

	// Determine the directory of the source file to place the output in the same directory
	outputDir := filepath.Dir(*srcFile)
	*output = filepath.Join(outputDir, *output)

	outputFile, err := os.Create(*output)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	data := struct {
		PackageName string
		Structs     []struct {
			TypeName string
			Fields   []field
		}
	}{
		PackageName: packageName,
	}

	for _, name := range structNames {
		data.Structs = append(data.Structs, struct {
			TypeName string
			Fields   []field
		}{
			TypeName: name,
			Fields:   structs[name],
		})
	}

	if err := t.Execute(outputFile, data); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

func getTypeName(expr ast.Expr) (string, string, bool) {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name, x.Name, false
	case *ast.SelectorExpr:
		typeName, _, _ := getTypeName(x.X)
		return typeName + "." + x.Sel.Name, x.Sel.Name, true
	case *ast.StarExpr:
		_, element, _ := getTypeName(x.X)
		return "*" + element, element, true
	case *ast.ArrayType:
		_, element, isStruct := getTypeName(x.Elt)
		return "[]" + element, element, isStruct
	default:
		return "", "", false
	}
}
