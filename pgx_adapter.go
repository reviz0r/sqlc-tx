package sqlc_tx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection interface {
	*pgx.Conn | *pgxpool.Conn | *pgxpool.Pool

	Begin(context.Context) (pgx.Tx, error)
}
