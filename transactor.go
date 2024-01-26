package sqlc_tx

import (
	"context"
	"errors"
)

type TransactorInterface[I any] interface {
	WithTx(context.Context, Func[I]) error
	WithoutTx(context.Context, Func[I]) error
}

type Transactor[C Connection, D any, Q Queries[Q], I any] struct {
	conn C

	querierConstructor Constructor[D, Q]
}

func New[C Connection, D any, Q Queries[Q], I any](conn C, c Constructor[D, Q]) *Transactor[C, D, Q, I] {
	return &Transactor[C, D, Q, I]{conn: conn, querierConstructor: c}
}

func (t *Transactor[C, D, Q, I]) WithTx(ctx context.Context, fn Func[I]) error {
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	txx := any(tx).(D)
	querier := t.querierConstructor(txx)
	querierInterface := any(querier).(I)

	err = fn(ctx, querierInterface)
	if err != nil {
		txErr := tx.Rollback(ctx)
		return errors.Join(err, txErr)
	}

	err = tx.Commit(ctx)
	if err != nil {
		txErr := tx.Rollback(ctx)
		return errors.Join(err, txErr)
	}

	return nil
}

func (t *Transactor[C, D, Q, I]) WithoutTx(ctx context.Context, fn Func[I]) error {
	conn := any(t.conn).(D)
	querier := t.querierConstructor(conn)
	querierInterface := any(querier).(I)
	return fn(ctx, querierInterface)
}
