package gen

import (
	"go/ast"
	"testing"

	"github.com/forbearing/golib/types/consts"
)

func TestIsServiceMethod1(t *testing.T) {
	fn1 := ServiceMethod1("u", "User", "CreateBefore", "model")
	fn2 := ServiceMethod2("u", "User", "ListBefore", "model")
	if !IsServiceMethod1(fn1) {
		t.Fatalf("expected IsServiceMethod1 to return true for ServiceMethod1-generated func")
	}
	if IsServiceMethod1(fn2) {
		t.Fatalf("expected IsServiceMethod1 to return false for non-matching func (ServiceMethod2)")
	}
}

func TestIsServiceMethod2(t *testing.T) {
	fn := ServiceMethod2("u", "User", "ListBefore", "model")
	fnNeg := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	if !IsServiceMethod2(fn) {
		t.Fatalf("expected IsServiceMethod2 to return true for ServiceMethod2-generated func")
	}
	if IsServiceMethod2(fnNeg) {
		t.Fatalf("expected IsServiceMethod2 to return false for non-matching func (ServiceMethod3)")
	}
}

func TestIsServiceMethod3(t *testing.T) {
	fn := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	fnNeg := ServiceMethod1("u", "User", "CreateBefore", "model")
	if !IsServiceMethod3(fn) {
		t.Fatalf("expected IsServiceMethod3 to return true for ServiceMethod3-generated func")
	}
	if IsServiceMethod3(fnNeg) {
		t.Fatalf("expected IsServiceMethod3 to return false for non-matching func (ServiceMethod1)")
	}
}

func TestIsServiceMethod4(t *testing.T) {
	fn := ServiceMethod4("u", "User", "Create", "model", "UserReq", "UserRsp")
	fnNeg := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	if !IsServiceMethod4(fn) {
		t.Fatalf("expected IsServiceMethod4 to return true for ServiceMethod4-generated func")
	}
	if IsServiceMethod4(fnNeg) {
		t.Fatalf("expected IsServiceMethod4 to return false for non-matching func (ServiceMethod3)")
	}
}

func TestIsServiceType(t *testing.T) {
	// Positive case: struct embeds service.Base[*model.User, *model.User, *model.User]
	gd := Types("model", "User", "User", "User", consts.PHASE_CREATE, false)
	if len(gd.Specs) == 0 {
		t.Fatalf("Types() returned no specs")
	}
	ts, ok := gd.Specs[0].(*ast.TypeSpec)
	if !ok {
		t.Fatalf("expected first spec to be *ast.TypeSpec")
	}
	if !IsServiceType(ts) {
		t.Fatalf("expected IsServiceType to return true for valid service.Base with three pointer type params")
	}

	// Negative case 1: wrong selector name (service.Other)
	neg1 := &ast.TypeSpec{
		Name: ast.NewIdent("userx"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: []*ast.Field{
				{Type: &ast.IndexListExpr{
					X: &ast.SelectorExpr{X: ast.NewIdent("service"), Sel: ast.NewIdent("Other")},
					Indices: []ast.Expr{
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
					},
				}},
			}},
		},
	}
	if IsServiceType(neg1) {
		t.Fatalf("expected IsServiceType to return false for non-Base selector")
	}

	// Negative case 2: one of the type params is not a pointer
	neg2 := &ast.TypeSpec{
		Name: ast.NewIdent("userx"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: []*ast.Field{
				{Type: &ast.IndexListExpr{
					X: &ast.SelectorExpr{X: ast.NewIdent("service"), Sel: ast.NewIdent("Base")},
					Indices: []ast.Expr{
						&ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}, // not a *T
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
					},
				}},
			}},
		},
	}
	if IsServiceType(neg2) {
		t.Fatalf("expected IsServiceType to return false when a type param is not a pointer")
	}

	// Negative case 3: incorrect number of type params (2 instead of 3)
	neg3 := &ast.TypeSpec{
		Name: ast.NewIdent("userx"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{List: []*ast.Field{
				{Type: &ast.IndexListExpr{
					X: &ast.SelectorExpr{X: ast.NewIdent("service"), Sel: ast.NewIdent("Base")},
					Indices: []ast.Expr{
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
						&ast.StarExpr{X: &ast.SelectorExpr{X: ast.NewIdent("model"), Sel: ast.NewIdent("User")}},
					},
				}},
			}},
		},
	}
	if IsServiceType(neg3) {
		t.Fatalf("expected IsServiceType to return false for wrong number of type params")
	}
}
