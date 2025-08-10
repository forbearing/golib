package gen

import (
	"testing"
)

func Test_IsServiceMethod1(t *testing.T) {
	fn1 := ServiceMethod1("u", "User", "CreateBefore", "model")
	fn2 := ServiceMethod2("u", "User", "ListBefore", "model")
	if !IsServiceMethod1(fn1) {
		t.Fatalf("expected IsServiceMethod1 to return true for ServiceMethod1-generated func")
	}
	if IsServiceMethod1(fn2) {
		t.Fatalf("expected IsServiceMethod1 to return false for non-matching func (ServiceMethod2)")
	}
}

func Test_IsServiceMethod2(t *testing.T) {
	fn := ServiceMethod2("u", "User", "ListBefore", "model")
	fnNeg := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	if !IsServiceMethod2(fn) {
		t.Fatalf("expected IsServiceMethod2 to return true for ServiceMethod2-generated func")
	}
	if IsServiceMethod2(fnNeg) {
		t.Fatalf("expected IsServiceMethod2 to return false for non-matching func (ServiceMethod3)")
	}
}

func Test_IsServiceMethod3(t *testing.T) {
	fn := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	fnNeg := ServiceMethod1("u", "User", "CreateBefore", "model")
	if !IsServiceMethod3(fn) {
		t.Fatalf("expected IsServiceMethod3 to return true for ServiceMethod3-generated func")
	}
	if IsServiceMethod3(fnNeg) {
		t.Fatalf("expected IsServiceMethod3 to return false for non-matching func (ServiceMethod1)")
	}
}

func Test_IsServiceMethod4(t *testing.T) {
	fn := ServiceMethod4("u", "User", "Create", "model", "UserReq", "UserRsp")
	fnNeg := ServiceMethod3("u", "User", "CreateManyBefore", "model")
	if !IsServiceMethod4(fn) {
		t.Fatalf("expected IsServiceMethod4 to return true for ServiceMethod4-generated func")
	}
	if IsServiceMethod4(fnNeg) {
		t.Fatalf("expected IsServiceMethod4 to return false for non-matching func (ServiceMethod3)")
	}
}
