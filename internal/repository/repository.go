package repository

import (
	"context"
	"github.com/allegro/bigcache"
	"github.com/jackc/pgx/v4"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/cache"
	"wb-l-zero/internal/repository/postgres"
)

type OrderPostgres interface {
	Save(ctx context.Context, order *entity.Order) error
	Get(ctx context.Context, orderUID string) (*entity.Order, error)
}

type OrderCache interface {
	GetCache(orderUID string) (*entity.Order, error)
	SaveCache(order *entity.Order) error
}

type Payment interface {
	Save(ctx context.Context, payment *entity.Payment) error
	Get(ctx context.Context, orderUID string) (*entity.Payment, error)
}

type Delivery interface {
	Save(ctx context.Context, delivery *entity.Delivery) error
	Get(ctx context.Context, orderUID string) (*entity.Delivery, error)
}

type Item interface {
	Save(ctx context.Context, item *entity.Item) error
	GetAllByOrderUID(ctx context.Context, orderUID string) ([]entity.Item, error)
}

type DBTransaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repository struct {
	OrderPostgres
	OrderCache
	Payment
	Delivery
	Item
	DBTransaction
}

func NewRepository(pgConn *pgx.Conn, cacheInp *bigcache.BigCache) *Repository {
	return &Repository{
		OrderPostgres: postgres.NewOrderPostgres(pgConn),
		OrderCache:    cache.NewOrderCache(cacheInp),
		Payment:       postgres.NewPaymentPostgres(pgConn),
		Delivery:      postgres.NewDeliveryPostgres(pgConn),
		Item:          postgres.NewItemPostgres(pgConn),
		DBTransaction: postgres.NewDBConn(pgConn),
	}
}
