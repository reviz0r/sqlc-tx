package sqlc_tx

import "context"

type Func[I any] func(context.Context, I) error

func Combine[I any](fns ...Func[I]) Func[I] {
	return func(ctx context.Context, i I) error {
		for _, fn := range fns {
			err := fn(ctx, i)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
