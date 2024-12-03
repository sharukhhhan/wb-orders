package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"

	"github.com/jackc/pgx/v4"
)

type PaymentPostgres struct {
	*pgx.Conn
}

func NewPaymentPostgres(conn *pgx.Conn) *PaymentPostgres {
	return &PaymentPostgres{Conn: conn}
}

func (p *PaymentPostgres) Save(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := p.Exec(ctx, query,
		payment.OrderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
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

func (r *PaymentPostgres) Get(ctx context.Context, orderUID string) (*entity.Payment, error) {
	query := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`
	var payment entity.Payment
	err := r.QueryRow(ctx, query, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}
