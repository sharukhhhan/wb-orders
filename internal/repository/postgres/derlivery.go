package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"
)

type DeliveryPostgres struct {
	*pgx.Conn
}

func NewDeliveryPostgres(conn *pgx.Conn) *DeliveryPostgres {
	return &DeliveryPostgres{Conn: conn}
}

func (r *DeliveryPostgres) Save(ctx context.Context, delivery *entity.Delivery) error {
	query := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.Exec(ctx, query,
		delivery.OrderUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repoerrors.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *DeliveryPostgres) Get(ctx context.Context, orderUID string) (*entity.Delivery, error) {
	query := `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`
	var delivery entity.Delivery
	err := r.QueryRow(ctx, query, orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &delivery, nil
}
