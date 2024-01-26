package sqlc_tx_test

import (
	"context"
	"log/slog"

	sqlc_tx "github.com/reviz0r/sqlc-tx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Example() {
	var pool *pgxpool.Pool

	s := sqlc_tx.New[*pgxpool.Pool, DBTX, *Queries, Querier](pool, New)

	var _ sqlc_tx.TransactorInterface[Querier] = s

	fn := sqlc_tx.Combine(
		func(ctx context.Context, querier Querier) error {
			_, err := querier.CreateOrganization(ctx, "hello world")
			return err

		},
		func(ctx context.Context, querier Querier) error {
			_, err := querier.GetOrganizationByID(ctx, pgtype.UUID{})
			return err

		},
	)

	err := s.WithTx(context.Background(), fn)
	if err != nil {
		slog.Error(err.Error())
	}
}

type Querier interface {
	CreateOrganization(ctx context.Context, name string) (pgtype.UUID, error)
	GetOrganizationByID(ctx context.Context, organizationID pgtype.UUID) (*Organization, error)
}

var _ Querier = (*Queries)(nil)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

type Organization struct {
	OrganizationID pgtype.UUID
	Name           string
}

const createOrganization = `-- name: CreateOrganization :one
insert into "organizations" ("name") values ($1) returning organization_id
`

func (q *Queries) CreateOrganization(ctx context.Context, name string) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createOrganization, name)
	var organization_id pgtype.UUID
	err := row.Scan(&organization_id)
	return organization_id, err
}

const getOrganizationByID = `-- name: GetOrganizationByID :one
select organization_id, name from "organizations" where organization_id = $1
`

func (q *Queries) GetOrganizationByID(ctx context.Context, organizationID pgtype.UUID) (*Organization, error) {
	row := q.db.QueryRow(ctx, getOrganizationByID, organizationID)
	var i Organization
	err := row.Scan(&i.OrganizationID, &i.Name)
	return &i, err
}
