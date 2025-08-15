package gen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/stoewer/go-strcase"
	goimports "golang.org/x/tools/imports"
	fumpt "mvdan.cc/gofumpt/format"
)

// FormatNode use go standard lib "go/format" to format ast.Node into code.
func FormatNode(node ast.Node, processImport ...bool) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formated, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	if len(processImport) > 0 && processImport[0] {
		result, err := goimports.Process("", formated, nil)
		if err != nil {
			return "", err
		}
		return string(result), nil
	}

	return string(formated), nil
}

// FormatNodeExtra use "https://github.com/mvdan/gofumpt" to format ast.Node into code.
func FormatNodeExtra(node ast.Node, processImport ...bool) (string, error) {
	return FormatNodeExtraWithFileSet(node, nil, processImport...)
}

// FormatNodeExtraWithFileSet 使用指定的 FileSet 格式化节点，保持注释位置
func FormatNodeExtraWithFileSet(node ast.Node, fset *token.FileSet, processImport ...bool) (string, error) {
	var buf bytes.Buffer
	// 如果没有提供 FileSet，创建新的
	if fset == nil {
		fset = token.NewFileSet()
	}

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formatted, err := fumpt.Source(buf.Bytes(), fumpt.Options{
		LangVersion: "",
		ExtraRules:  true,
	})

	if len(processImport) > 0 && processImport[0] {
		var result []byte
		if result, err = goimports.Process("", formatted, nil); err != nil {
			return "", err
		}
		return string(result), nil
	}

	return string(formatted), err
}

func MethodAddComments(code string, modelName string) string {
	for _, method := range Methods {
		str := strings.ReplaceAll(strcase.SnakeCase(method), "_", " ")
		// Add comment after log.Info
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
