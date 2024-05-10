package main

import (
	"flag"
	"fmt"
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
        {{- if .IsPointer }}
			if src.{{ .Name }} != nil {
				if dst.{{ .Name }} == nil {
					dst.{{ .Name }} = new({{ .TypeElement }})
				}
				Merge{{ .TypeElement }}(dst.{{ .Name }}, src.{{ .Name }})
			}
        {{- else }}
    		Merge{{ .TypeElement }}(&dst.{{ .Name }}, &src.{{ .Name }})
        {{- end }}
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
	IsPointer   bool
}

type templateData struct {
	PackageName string
	Structs     []structData
}

type structData struct {
	TypeName string
	Fields   []field
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

	t, err := parseTemplate()
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return
	}

	outputFile, err := generateOutputFile(*srcFile, *output)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	data := TransformAstToTemplateData(node)

	if err := t.Execute(outputFile, data); err != nil {
		fmt.Println("Error executing template:", err)
	}
}

func generateOutputFile(src string, out string) (*os.File, error) {
	// Determine the directory of the source file to place the output in the same directory
	outputDir := filepath.Dir(src)
	output := filepath.Join(outputDir, out)
	outputFile, err := os.Create(output)
	return outputFile, err
}

func parseTemplate() (*template.Template, error) {
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

	return template.New("func").Funcs(funcMap).Parse(funcTemplate)
}
