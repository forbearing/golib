package service

import (
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

func init() {
	service.Register[*star](consts.PHASE_CREATE)
}

type star struct {
	service.Base[*model.Star, *model.Star, *model.Star]
}

func (s *star) CreateBefore(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star create before")
	return nil
}

func (s *star) CreateAfter(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star create after")
	return nil
}

func (s *star) DeleteBefore(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star delete before")
	return nil
}

func (s *star) DeleteAfter(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star delete after")
	return nil
}

func (s *star) UpdateBefore(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star update before")
	return nil
}

func (s *star) UpdateAfter(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star update after")
	return nil
}

func (s *star) PatchBefore(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star update partial before")
	return nil
}

func (s *star) PatchAfter(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star update partial after")
	return nil
}

func (s *star) ListBefore(ctx *types.ServiceContext, stars *[]*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star list before")
	return nil
}

func (s *star) ListAfter(ctx *types.ServiceContext, stars *[]*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star list after")
	return nil
}

func (s *star) GetBefore(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star get before")
	return nil
}

func (s *star) GetAfter(ctx *types.ServiceContext, star *model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star get after")
	return nil
}

func (s *star) CreateManyBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch create before")
	return nil
}

func (s *star) CreateManyAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch create after")
	return nil
}

func (s *star) DeleteManyBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch delete before")
	return nil
}

func (s *star) DeleteManyAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch delete after")
	return nil
}

func (s *star) UpdateManyBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch update before")
	return nil
}

func (s *star) UpdateManyAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch update after")
	return nil
}

func (s *star) PatchManyBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch update partial before")
	return nil
}

func (s *star) PatchManyAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("star batch update partial after")
	return nil
}
