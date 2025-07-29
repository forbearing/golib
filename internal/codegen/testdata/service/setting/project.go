package service_setting

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_setting "nebula/model/setting"
)

func init() {
	service.Register[*project]()
}

// project implements the types.Service[*model_setting.Project] interface.
type project struct {
	service.Base[*model_setting.Project]
}

func (p *project) CreateBefore(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) CreateAfter(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) DeleteBefore(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) DeleteAfter(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) UpdateBefore(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) UpdateAfter(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) UpdatePartialBefore(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) UpdatePartialAfter(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) ListBefore(ctx *types.ServiceContext, projects *[]*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) ListAfter(ctx *types.ServiceContext, projects *[]*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) GetBefore(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) GetAfter(ctx *types.ServiceContext, project *model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchCreateBefore(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchCreateAfter(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchDeleteBefore(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchDeleteAfter(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchUpdateBefore(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchUpdateAfter(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchUpdatePartialBefore(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (p *project) BatchUpdatePartialAfter(ctx *types.ServiceContext, projects ...*model_setting.Project) error {
	log := p.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("project batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
