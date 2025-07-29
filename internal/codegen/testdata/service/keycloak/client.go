package service_keycloak

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_keycloak "nebula/model/keycloak"
)

func init() {
	service.Register[*client]()
}

// client implements the types.Service[*model_keycloak.Client] interface.
type client struct {
	service.Base[*model_keycloak.Client]
}

func (c *client) CreateBefore(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) CreateAfter(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) DeleteBefore(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) DeleteAfter(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) UpdateBefore(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) UpdateAfter(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) UpdatePartialBefore(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) UpdatePartialAfter(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) ListBefore(ctx *types.ServiceContext, clients *[]*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) ListAfter(ctx *types.ServiceContext, clients *[]*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) GetBefore(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) GetAfter(ctx *types.ServiceContext, client *model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchCreateBefore(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchCreateAfter(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchDeleteBefore(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchDeleteAfter(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchUpdateBefore(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchUpdateAfter(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchUpdatePartialBefore(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (c *client) BatchUpdatePartialAfter(ctx *types.ServiceContext, clients ...*model_keycloak.Client) error {
	log := c.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("client batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
