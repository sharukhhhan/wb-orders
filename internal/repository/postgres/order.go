package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"
)

type OrderPostgres struct {
	*pgx.Conn
}

func NewOrderPostgres(conn *pgx.Conn) *OrderPostgres {
	return &OrderPostgres{conn}
}

func (p *OrderPostgres) Save(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, customer_id, 
			delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := p.Exec(ctx, query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OOFShard,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return repoerrors.ErrAlreadyExists
		}
		return err
	}

	return nil
}

func (r *OrderPostgres) Get(ctx context.Context, orderUID string) (*entity.Order, error) {
	query := `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`
	var order entity.Order
	err := r.QueryRow(ctx, query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SMID,
		&order.DateCreated,
		&order.OOFShard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}
