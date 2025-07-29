package gen

import (
	"bytes"
	"go/format"
	"go/token"
	"testing"
)

func Test_expr_log_info(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		str  string
		want string
	}{
		{
			str:  `"hello world"`,
			want: `log.Info("hello world")`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			res := expr_log_info(tt.str)
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("expr_log_info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_return_one(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		str  string
		want string
	}{
		{
			name: "return error",
			str:  `error`,
			want: `return error`,
		},
		{
			name: "return nil",
			str:  `nil`,
			want: `return nil`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := returns(tt.str)
			buf.Reset()
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("return_one() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assign_with_service_context(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modelVarName string
		want         string
	}{
		{
			name:         "u",
			modelVarName: `u`,
			want:         `log := u.WithServiceContext(ctx, ctx.GetPhase())`,
		},
		{
			name:         "g",
			modelVarName: `g`,
			want:         `log := g.WithServiceContext(ctx, ctx.GetPhase())`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := assign_with_service_context(tt.modelVarName)
			buf.Reset()
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("assign_with_service_context() = %v, want %v", got, tt.want)
			}
		})
	}
}
