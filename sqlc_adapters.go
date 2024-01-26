package sqlc_tx

import "github.com/jackc/pgx/v5"

type Queries[Q any] interface {
	WithTx(tx pgx.Tx) Q
}

type Constructor[D any, Q Queries[Q]] func(D) Q
