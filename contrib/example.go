package contrib

import "context"

func ok() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx := context.WithValue(ctx, "key", i)
		_ = ctx
	}
}

func notOk() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i) // "nested context in loop"
		_ = ctx
	}
}
