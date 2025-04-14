package service

import (
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

func init() {
	service.Register[*star]()
}

type star struct {
	service.Base[*model.Star]
}

func (s *star) CreateBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_CREATE_BEFORE)
	log.Info("star create before")
	return nil
}

func (s *star) CreateAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_CREATE_AFTER)
	log.Info("star create after")
	return nil
}

func (s *star) DeleteBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_DELETE_BEFORE)
	log.Info("star delete before")
	return nil
}

func (s *star) DeleteAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	log.Info("star delete after")
	return nil
}

func (s *star) UpdateBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_UPDATE_BEFORE)
	log.Info("star update before")
	return nil
}

func (s *star) UpdateAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_UPDATE_AFTER)
	log.Info("star update after")
	return nil
}

func (s *star) UpdatePartialBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_BEFORE)
	log.Info("star update partial before")
	return nil
}

func (s *star) UpdatePartialAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_AFTER)
	log.Info("star update partial after")
	return nil
}

func (s *star) ListBefore(ctx *types.ServiceContext, stars *[]*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_LIST_BEFORE)
	log.Info("star list before")
	return nil
}

func (s *star) ListAfter(ctx *types.ServiceContext, stars *[]*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_LIST_AFTER)
	log.Info("star list after")
	return nil
}

func (s *star) GetBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_GET_BEFORE)
	log.Info("star get before")
	return nil
}

func (s *star) GetAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_GET_AFTER)
	log.Info("star get after")
	return nil
}

func (s *star) BatchCreateBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_BEFORE)
	log.Info("star batch create before")
	return nil
}

func (s *star) BatchCreateAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_AFTER)
	log.Info("star batch create after")
	return nil
}

func (s *star) BatchDeleteBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_BEFORE)
	log.Info("star batch delete before")
	return nil
}

func (s *star) BatchDeleteAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_AFTER)
	log.Info("star batch delete after")
	return nil
}

func (s *star) BatchUpdateBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_BEFORE)
	log.Info("star batch update before")
	return nil
}

func (s *star) BatchUpdateAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_AFTER)
	log.Info("star batch update after")
	return nil
}

func (s *star) BatchUpdatePartialBefore(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_BEFORE)
	log.Info("star batch update partial before")
	return nil
}

func (s *star) BatchUpdatePartialAfter(ctx *types.ServiceContext, stars ...*model.Star) error {
	log := s.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_AFTER)
	log.Info("star batch update partial after")
	return nil
}
