package service

import (
	"context"
	"errors"
	"fmt"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository"
	"wb-l-zero/internal/repository/repoerrors"
)

type OrderService struct {
	orderRepoPostgres repository.OrderPostgres
	orderRepoCache    repository.OrderCache
	deliveryRepo      repository.Delivery
	paymentRepo       repository.Payment
	itemRepo          repository.Item
	t                 repository.DBTransaction
}

func NewOrderService(orderRepo repository.OrderPostgres, oderRepoCache repository.OrderCache, deliveryRepo repository.Delivery, paymentRepo repository.Payment, itemRepo repository.Item, t repository.DBTransaction) *OrderService {
	return &OrderService{orderRepoPostgres: orderRepo, orderRepoCache: oderRepoCache, deliveryRepo: deliveryRepo, paymentRepo: paymentRepo, itemRepo: itemRepo, t: t}
}

func (s *OrderService) Create(ctx context.Context, order *entity.Order) error {

	tx, err := s.t.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	err = s.orderRepoPostgres.Save(ctx, order)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExists) {
			err = ErrOrderAlreadyExists
		}
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to save order into postgres: %w", err)
	}

	order.Delivery.OrderUID = order.OrderUID
	err = s.deliveryRepo.Save(ctx, &order.Delivery)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to save delivery: %w", err)
	}

	order.Payment.OrderUID = order.OrderUID
	err = s.paymentRepo.Save(ctx, &order.Payment)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to save payment: %w", err)
	}

	for _, item := range order.Items {
		item.OrderUID = order.OrderUID
		err = s.itemRepo.Save(ctx, &item)
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed to save items: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	err = s.orderRepoCache.SaveCache(order)
	if err != nil {
		return fmt.Errorf("failed to save order into cache: %w", err)
	}

	return nil
}

func (s *OrderService) GetOrderDetails(ctx context.Context, orderUID string) (*entity.Order, error) {

	order, err := s.orderRepoCache.GetCache(orderUID)
	if err != nil {
		if !errors.Is(err, repoerrors.ErrNotFound) {
			return nil, fmt.Errorf("failed to get order from cache: %w", err)
		}
	}
	if order != nil {
		return order, nil
	}

	order, err = s.orderRepoPostgres.Get(ctx, orderUID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			err = ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order from postgres: %w", err)
	}

	delivery, err := s.deliveryRepo.Get(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	order.Delivery = *delivery

	payment, err := s.paymentRepo.Get(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	order.Payment = *payment

	items, err := s.itemRepo.GetAllByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	order.Items = items

	return order, nil
}
