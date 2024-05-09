package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// helper function to parse source to *ast.File
func parseSource(src string) *ast.File {
	node, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)
	if err != nil {
		panic("Failed to parse source: " + err.Error())
	}
	return node
}

func TestTransformAstToTemplateData(t *testing.T) {
	src := `
		package test

		type Features struct {
			EnableLogging bool
			MaxRetries    int
		}
		
		type Client struct {
			Host string
			Port int
		}
		
		type Config struct {
			APIKey   string
			Timeout  int
			Features *Features
			Client   Client
			Bar      string
		}
	`

	node := parseSource(src)

	result := TransformAstToTemplateData(node)

	if result.PackageName != "test" {
		t.Errorf("Expected test, got %s", result.PackageName)
	}
	for _, structData := range result.Structs {
		if structData.TypeName == "Config" {
			features := structData.Fields[2]
			if features.IsPointer != true {
				t.Errorf("Expected true, got %t", features.IsPointer)
			}
			if features.IsStruct != true {
				t.Errorf("Expected true, got %t", features.IsStruct)
			}
		}
	}
}
