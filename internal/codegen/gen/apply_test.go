package gen

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/forbearing/golib/dsl"
	"github.com/kr/pretty"
)

func TestApplyServiceFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		code   string
		action *dsl.Action
		want   string
	}{
		{
			name: "user_create_with_payload_result",
			code: dataServiceUserCreate,
			action: &dsl.Action{
				Enabled: true,
				Payload: "UserReq",
				Result:  "UserRsp",
			},
			want: `package service

import (
	"helloworld/model"

	"github.com/forbearing/golib/types"
)

func (u *user) Create(ctx *types.ServiceContext, req *model.UserReq) (rsp *model.UserRsp, err error) {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create")
	return rsp, nil
}

func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create after")
	return nil
}
`,
		},
		{
			name: "user_create_no_payload_result",
			code: dataServiceUserCreate,
			action: &dsl.Action{
				Enabled: true,
				Payload: "User",
				Result:  "User",
			},
			want: `package service

import (
	"helloworld/model"

	"github.com/forbearing/golib/types"
)

func (u *user) Create(ctx *types.ServiceContext, req *model.User) (rsp *model.User, err error) {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create")
	return rsp, nil
}

func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create after")
	return nil
}
`,
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
			ApplyServiceFile(file, tt.action)
			got, err := FormatNodeExtra(file)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}
