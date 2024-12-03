package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/repository/repoerrors"
)

type ItemPostgres struct {
	*pgx.Conn
}

func NewItemPostgres(conn *pgx.Conn) *ItemPostgres {
	return &ItemPostgres{conn}
}

func (p *ItemPostgres) Save(ctx context.Context, item *entity.Item) error {
	query := `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := p.Exec(ctx, query,
		item.OrderUID,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.RID,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status,
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

func (r *ItemPostgres) GetAllByOrderUID(ctx context.Context, orderUID string) ([]entity.Item, error) {
	query := `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`
	rows, err := r.Query(ctx, query, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.Item
	for rows.Next() {
		var item entity.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
