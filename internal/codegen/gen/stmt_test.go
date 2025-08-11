package gen

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"testing"

	"github.com/forbearing/golib/types/consts"
)

func TestStmtLogInfo(t *testing.T) {
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
			res := StmtLogInfo(tt.str)
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("StmtLogInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReturns(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		exprs []ast.Expr
		want  string
	}{
		{
			name:  "return error",
			exprs: []ast.Expr{ast.NewIdent("error")},
			want:  `return error`,
		},
		{
			name:  "return nil",
			exprs: []ast.Expr{ast.NewIdent("nil")},
			want:  `return nil`,
		},
		{
			name: "return &model.User{}, nil",
			exprs: []ast.Expr{
				&ast.UnaryExpr{
					Op: token.AND,
					X: &ast.CompositeLit{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("model"),
							Sel: ast.NewIdent("User"),
						},
					},
				},
				ast.NewIdent("nil"),
			},
			want: "return &model.User{}, nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Returns(tt.exprs...)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("Returns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStmtLogWithServiceContext(t *testing.T) {
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
			res := StmtLogWithServiceContext(tt.modelVarName)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("StmtLogWithServiceContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStmtRouterRegister(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modelPkgName string
		modelName    string
		reqName      string
		respName     string
		endpoint     string
		verb         string
		want         string
	}{
		{
			name:         "test1",
			modelPkgName: "model",
			modelName:    "Group",
			reqName:      "Group",
			respName:     "Group",
			endpoint:     "group",
			verb:         "Create",
			want:         `router.Register[*model.Group, *model.Group, *model.Group](router.API(), "group", consts.Create)`,
		},
		{
			name:         "test2",
			modelPkgName: "pkgmodel",
			modelName:    "Group",
			reqName:      "GroupRequest",
			respName:     "GroupResponse",
			endpoint:     "group2",
			verb:         "Update",
			want:         `router.Register[*pkgmodel.Group, *pkgmodel.GroupRequest, *pkgmodel.GroupResponse](router.API(), "group2", consts.Update)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := StmtRouterRegister(tt.modelPkgName, tt.modelName, tt.reqName, tt.respName, tt.endpoint, tt.verb)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("StmtRouterRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStmtServiceRegister(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		structName string
		want       string
		phase      consts.Phase
	}{
		{
			name:       "test1",
			structName: "user",
			phase:      consts.PHASE_CREATE,
			want:       `service.Register[*user](consts.PHASE_CREATE)`,
		},
		{
			name:       "test2",
			structName: "group",
			phase:      consts.PHASE_UPDATE,
			want:       `service.Register[*group](consts.PHASE_UPDATE)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := StmtServiceRegister(tt.structName, tt.phase)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("StmtServiceRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}
