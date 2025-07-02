// internal/codegen/generator/template.go
package generator

const ServiceTemplate = `package {{.PackageName}}

import (
	"{{.ModelPackage}}"

	"{{.FrameworkPath}}/service"
	"{{.FrameworkPath}}/types"
)

func init() {
	service.Register[*{{.ServiceName}}]()
}

// {{.ServiceName}} implements the types.Service[*model.{{.ModelName}}] interface.
type {{.ServiceName}} struct {
	service.Base[*model.{{.ModelName}}]
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) CreateBefore(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) CreateAfter(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) DeleteBefore(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) DeleteAfter(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) UpdateBefore(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) UpdateAfter(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) UpdatePartialBefore(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) UpdatePartialAfter(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) ListBefore(ctx *types.ServiceContext, {{.ModelVariable}} *[]*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) ListAfter(ctx *types.ServiceContext, {{.ModelVariable}} *[]*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) GetBefore(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) GetAfter(ctx *types.ServiceContext, {{.ModelVariable}} *model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchCreateBefore(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchCreateAfter(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchDeleteBefore(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchDeleteAfter(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchUpdateBefore(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchUpdateAfter(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchUpdatePartialBefore(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func ({{.ServiceName | firstChar}} *{{.ServiceName}}) BatchUpdatePartialAfter(ctx *types.ServiceContext, {{.ModelVariable}} ...*model.{{.ModelName}}) error {
	log := {{.ServiceName | firstChar}}.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("{{.ServiceName}} batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
`
