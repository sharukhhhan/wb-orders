package service

import (
	"context"
	"testing"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderDetails_FromCache(t *testing.T) {
	mockOrder := &entity.Order{OrderUID: "12345"}
	cacheRepo := &MockOrderRepoCache{
		GetFunc: func(orderUID string) (*entity.Order, error) {
			if orderUID == "12345" {
				return mockOrder, nil
			}
			return nil, repoerrors.ErrNotFound
		},
	}
	postgresRepo := &MockOrderRepoPostgres{}

	service := &OrderService{
		orderRepoCache:    cacheRepo,
		orderRepoPostgres: postgresRepo,
	}

	order, err := service.GetOrderDetails(context.Background(), "12345")
	assert.NoError(t, err)
	assert.Equal(t, mockOrder, order)
}

func TestGetOrderDetails_FromPostgres(t *testing.T) {
	mockOrder := &entity.Order{OrderUID: "12345"}
	cacheRepo := &MockOrderRepoCache{
		GetFunc: func(orderUID string) (*entity.Order, error) {
			return nil, repoerrors.ErrNotFound // Нет данных в кэше
		},
		SaveFunc: func(order *entity.Order) error {
			assert.Equal(t, mockOrder, order) // Проверка, что данные кэшируются
			return nil
		},
	}
	postgresRepo := &MockOrderRepoPostgres{
		GetFunc: func(ctx context.Context, orderUID string) (*entity.Order, error) {
			if orderUID == "12345" {
				return mockOrder, nil
			}
			return nil, repoerrors.ErrNotFound
		},
	}

	service := &OrderService{
		orderRepoCache:    cacheRepo,
		orderRepoPostgres: postgresRepo,
		// Убедимся, что другие репозитории инициализированы, даже если они не используются в этом тесте
		deliveryRepo: &MockDeliveryRepo{
			GetFunc: func(ctx context.Context, orderUID string) (*entity.Delivery, error) {
				return &entity.Delivery{}, nil
			},
		},
		paymentRepo: &MockPaymentRepo{
			GetFunc: func(ctx context.Context, orderUID string) (*entity.Payment, error) {
				return &entity.Payment{}, nil
			},
		},
		itemRepo: &MockItemRepo{
			GetAllByOrderUIDFunc: func(ctx context.Context, orderUID string) ([]entity.Item, error) {
				return []entity.Item{}, nil
			},
		},
	}

	order, err := service.GetOrderDetails(context.Background(), "12345")
	assert.NoError(t, err)
	assert.Equal(t, mockOrder, order)
}

func TestGetOrderDetails_OrderNotFound(t *testing.T) {
	cacheRepo := &MockOrderRepoCache{
		GetFunc: func(orderUID string) (*entity.Order, error) {
			return nil, repoerrors.ErrNotFound // Нет данных в кэше
		},
	}
	postgresRepo := &MockOrderRepoPostgres{
		GetFunc: func(ctx context.Context, orderUID string) (*entity.Order, error) {
			return nil, repoerrors.ErrNotFound // Нет данных в базе
		},
	}

	service := &OrderService{
		orderRepoCache:    cacheRepo,
		orderRepoPostgres: postgresRepo,
	}

	order, err := service.GetOrderDetails(context.Background(), "12345")
	assert.Nil(t, order)
	assert.ErrorIs(t, err, ErrOrderNotFound)
}
