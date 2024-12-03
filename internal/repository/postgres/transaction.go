package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DBConn struct {
	db *pgx.Conn
}

func NewDBConn(db *pgx.Conn) *DBConn {
	return &DBConn{db: db}
}

func (dbc *DBConn) Begin(ctx context.Context) (pgx.Tx, error) {
	return dbc.db.Begin(ctx)
}
