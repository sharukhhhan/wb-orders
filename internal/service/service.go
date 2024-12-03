package service

import (
	"context"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository"
)

type Order interface {
	Create(ctx context.Context, order *entity.Order) error
	GetOrderDetails(ctx context.Context, orderUID string) (*entity.Order, error)
}

type Service struct {
	Order
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repository.OrderPostgres, repository.OrderCache, repository.Delivery, repository.Payment, repository.Item, repository.DBTransaction),
	}
}
