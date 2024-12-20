package model

import pkgmodel "github.com/forbearing/golib/model"

type OrderStatus string

const (
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusShapped   OrderStatus = "shapped"
)

func init() {
	pkgmodel.Register[*Order]()
}

type Order struct {
	pkgmodel.Base

	UserID string      `json:"user_id"`
	Status OrderStatus `json:"status"`
}
