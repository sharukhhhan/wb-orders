package service

import (
	"context"
	"wb-l-zero/internal/entity"
)

type MockOrderRepoCache struct {
	GetFunc    func(orderUID string) (*entity.Order, error)
	SaveFunc   func(order *entity.Order) error
	DeleteFunc func(orderUID string)
}

func (m *MockOrderRepoCache) GetCache(orderUID string) (*entity.Order, error) {
	return m.GetFunc(orderUID)
}

func (m *MockOrderRepoCache) SaveCache(order *entity.Order) error {
	return m.SaveFunc(order)
}

type MockOrderRepoPostgres struct {
	GetFunc  func(ctx context.Context, orderUID string) (*entity.Order, error)
	SaveFunc func(ctx context.Context, order *entity.Order) error
}

func (m *MockOrderRepoPostgres) Get(ctx context.Context, orderUID string) (*entity.Order, error) {
	return m.GetFunc(ctx, orderUID)
}

func (m *MockOrderRepoPostgres) Save(ctx context.Context, order *entity.Order) error {
	return m.SaveFunc(ctx, order)
}

type MockDeliveryRepo struct {
	GetFunc  func(ctx context.Context, orderUID string) (*entity.Delivery, error)
	SaveFunc func(ctx context.Context, delivery *entity.Delivery) error
}

func (m *MockDeliveryRepo) Get(ctx context.Context, orderUID string) (*entity.Delivery, error) {
	return m.GetFunc(ctx, orderUID)
}

func (m *MockDeliveryRepo) Save(ctx context.Context, delivery *entity.Delivery) error {
	return m.SaveFunc(ctx, delivery)
}

type MockPaymentRepo struct {
	GetFunc  func(ctx context.Context, orderUID string) (*entity.Payment, error)
	SaveFunc func(ctx context.Context, payment *entity.Payment) error
}

func (m *MockPaymentRepo) Get(ctx context.Context, orderUID string) (*entity.Payment, error) {
	return m.GetFunc(ctx, orderUID)
}

func (m *MockPaymentRepo) Save(ctx context.Context, payment *entity.Payment) error {
	return m.SaveFunc(ctx, payment)
}

type MockItemRepo struct {
	GetAllByOrderUIDFunc func(ctx context.Context, orderUID string) ([]entity.Item, error)
	SaveFunc             func(ctx context.Context, item *entity.Item) error
}

func (m *MockItemRepo) GetAllByOrderUID(ctx context.Context, orderUID string) ([]entity.Item, error) {
	return m.GetAllByOrderUIDFunc(ctx, orderUID)
}

func (m *MockItemRepo) Save(ctx context.Context, item *entity.Item) error {
	return m.SaveFunc(ctx, item)
}
