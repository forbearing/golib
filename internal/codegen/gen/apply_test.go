package gen_test

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/internal/codegen/gen"
)

func TestApplyServiceFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code   string
		action *dsl.Action
	}{
		{
			name: "user_create",
			code: dataServiceUserCreate,
			action: &dsl.Action{
				Enabled: true,
				Payload: "User",
				Result:  "User",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "", tt.code, parser.ParseComments)
			if err != nil {
				t.Error(err)
				return
			}
			gen.ApplyServiceFile(file, tt.action)
		})
	}
}
