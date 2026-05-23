package fetcher

import "context"

type Result[T any] struct {
	Data  T
	Error error
}

func Fetch[T any](ctx context.Context, fetchFunc func(context.Context) Result[T]) <-chan Result[T] {
	resultChan := make(chan Result[T], 1)

	go func() {
		defer close(resultChan)
		resultChan <- fetchFunc(ctx)
	}()

	return resultChan
}
