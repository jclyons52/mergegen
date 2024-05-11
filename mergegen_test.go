package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
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
			Timeout  *int
			Features *Features
			Client   Client
			values   []int
			CreatedAt time.Time
		}
	`

	node := parseSource(src)
	result := TransformAstToTemplateData(node)

	assert.Equal(t, "test", result.PackageName, "PackageName should match 'test'")

	for _, structData := range result.Structs {
		if structData.TypeName == "Config" {
			assert.False(t, structData.Fields[0].IsPointer, "APIKey should not be a pointer")
			assert.False(t, structData.Fields[0].IsStruct, "APIKey should not be a struct")
			assert.True(t, structData.Fields[1].IsPointer, "Timeout should be a pointer")
			assert.False(t, structData.Fields[1].IsStruct, "Timeout should not be a struct")
			assert.True(t, structData.Fields[2].IsPointer, "Features should be a pointer")
			assert.True(t, structData.Fields[2].IsStruct, "Features should be a struct")
			assert.False(t, structData.Fields[3].IsPointer, "Client should not be a pointer")
			assert.True(t, structData.Fields[3].IsStruct, "Client should be a struct")
			assert.False(t, structData.Fields[4].IsPointer, "values should not be a pointer")
			assert.False(t, structData.Fields[4].IsStruct, "values should not be a struct")
			assert.True(t, structData.Fields[5].IsStruct, "CreatedAt should be a struct")
			assert.True(t, structData.Fields[5].IsExternal, "CreatedAt should be an external type")
		}
	}
}
