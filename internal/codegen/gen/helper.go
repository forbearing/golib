package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/stoewer/go-strcase"
	fumpt "mvdan.cc/gofumpt/format"
)

// FormatNode use go standard lib "go/format" to format ast.Node into code.
func FormatNode(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formated, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}
	return string(formated), nil
}

// FormatNodeExtra use "https://github.com/mvdan/gofumpt" to format ast.Node into code.
func FormatNodeExtra(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formatted, err := fumpt.Source(buf.Bytes(), fumpt.Options{
		LangVersion: "",
		ExtraRules:  true,
	})
	return string(formatted), err
}

func MethodAddComments(code string, modelName string) string {
	for _, method := range Methods {
		str := strings.ReplaceAll(strcase.SnakeCase(method), "_", " ")
		// 在 log.Info 之后添加注释
		searchStr := fmt.Sprintf(`log.Info("%s %s")`, strings.ToLower(modelName), str)
		replaceStr := fmt.Sprintf(`log.Info("%s %s")
	// =============================
	// Add your business logic here.
	// =============================
`, strings.ToLower(modelName), str)

		code = strings.ReplaceAll(code, searchStr, replaceStr)
	}

	return code
}

// formatAndImports formats code use gofumpt and processes imports.
func formatAndImports(f *ast.File) (string, error) {
	formatted, err := FormatNodeExtra(f)
	if err != nil {
		return "", err
	}

	// result, err := goimports.Process("", []byte(formatted), nil)
	// if err != nil {
	// 	return "", err
	// }

	return string(formatted), nil
}
