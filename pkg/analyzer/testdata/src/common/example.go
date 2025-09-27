package common

import (
	"context"
	"fmt"
	"testing"
)

func example() {
	ctx := context.Background()
	ctx2 := context.Background()

	for i := 0; i < 10; i++ {
		ctx := context.WithValue(ctx, "key", i)
		ctx = context.WithValue(ctx, "other", "val")
	}

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i) // want "nested context in loop"
		ctx = context.WithValue(ctx, "other", "val")
	}

	for i := 0; i < 10; i++ {
		ctx = context.Background()
		// This is fine because the first action was to re-assign to a known empty context.
		ctx = context.WithValue(ctx, "key", i)
		ctx = context.WithValue(ctx, "other", "val")
		// Not OK because this var has NOT been re-assigned to a known empty context.
		ctx2 = context.WithValue(ctx2, "key", i) // want "nested context in loop"
	}

	for item := range []string{"one", "two", "three"} {
		ctx = wrapContext(ctx) // want "nested context in loop"
		ctx := context.WithValue(ctx, "key", item)
		ctx = wrapContext(ctx)
	}

	for {
		ctx = wrapContext(ctx) // want "nested context in loop"
		break
	}

	// not fooled by shadowing in nested blocks
	for {
		err := doSomething()
		if err != nil {
			ctx := wrapContext(ctx)
			ctx = wrapContext(ctx)
		}

		switch err {
		case nil:
			ctx := wrapContext(ctx)
			ctx = wrapContext(ctx)
		default:
			ctx := wrapContext(ctx)
			ctx = wrapContext(ctx)
		}

		{
			ctx := wrapContext(ctx)
			ctx = wrapContext(ctx)
		}

		select {
		case <-ctx.Done():
			ctx := wrapContext(ctx)
			ctx = wrapContext(ctx)
		default:
		}

		ctx = wrapContext(ctx) // want "nested context in loop"

		break
	}

	// detects contexts wrapped in function literals (this is risky as function literals can be called multiple times)
	_ = func() {
		ctx = wrapContext(ctx) // want "nested context in function literal"
	}

	// this is fine because context.Background is an empty context
	_ = func() {
		ctx = context.Background()
	}

	_ = func(tb testing.TB) {
		ctx = tb.Context()
	}

	_ = func(b *testing.B) {
		ctx = b.Context()
	}

	_ = func(t *testing.T) {
		ctx = t.Context()
	}

	// this is fine because the context is created in the loop
	for {
		if ctx := context.Background(); doSomething() != nil {
			ctx = wrapContext(ctx)
		}
	}

	for {
		ctx2 := context.Background()
		ctx = wrapContext(ctx) // want "nested context in loop"
		if doSomething() != nil {
			ctx2 = wrapContext(ctx2)
		}
	}
}

func wrapContext(ctx context.Context) context.Context {
	return context.WithoutCancel(ctx)
}

func doSomething() error {
	return nil
}

// storing contexts in a struct isn't recommended, but local copies of a non-pointer struct should act like local copies of a context.
func inStructs(ctx context.Context) {
	for i := 0; i < 10; i++ {
		c := struct{ Ctx context.Context }{ctx}
		c.Ctx = context.WithValue(c.Ctx, "key", i)
		c.Ctx = context.WithValue(c.Ctx, "other", "val")
	}

	for i := 0; i < 10; i++ {
		c := []struct{ Ctx context.Context }{{ctx}}
		c[0].Ctx = context.WithValue(c[0].Ctx, "key", i)
		c[0].Ctx = context.WithValue(c[0].Ctx, "other", "val")
	}

	c := struct{ Ctx context.Context }{ctx}
	for i := 0; i < 10; i++ {
		c := c
		c.Ctx = context.WithValue(c.Ctx, "key", i)
		c.Ctx = context.WithValue(c.Ctx, "other", "val")
	}

	pc := &struct{ Ctx context.Context }{ctx}
	for i := 0; i < 10; i++ {
		c := pc
		c.Ctx = context.WithValue(c.Ctx, "key", i) // want "nested context in loop"
		c.Ctx = context.WithValue(c.Ctx, "other", "val")
	}

	r := []struct{ Ctx context.Context }{{ctx}}
	for i := 0; i < 10; i++ {
		r[0].Ctx = context.WithValue(r[0].Ctx, "key", i) // want "nested context in loop"
		r[0].Ctx = context.WithValue(r[0].Ctx, "other", "val")
	}

	rp := []*struct{ Ctx context.Context }{{ctx}}
	for i := 0; i < 10; i++ {
		rp[0].Ctx = context.WithValue(rp[0].Ctx, "key", i) // want "nested context in loop"
		rp[0].Ctx = context.WithValue(rp[0].Ctx, "other", "val")
	}
}

func inVariousNestedBlocks(ctx context.Context) {
	for {
		err := doSomething()
		if err != nil {
			ctx = wrapContext(ctx) // want "nested context in loop"
		}

		break
	}

	for {
		err := doSomething()
		if err != nil {
			if true {
				ctx = wrapContext(ctx) // want "nested context in loop"
			}
		}

		break
	}

	for {
		err := doSomething()
		switch err {
		case nil:
			ctx = wrapContext(ctx) // want "nested context in loop"
		}

		break
	}

	for {
		err := doSomething()
		switch err {
		default:
			ctx = wrapContext(ctx) // want "nested context in loop"
		}

		break
	}

	for {
		ctx := wrapContext(ctx)

		err := doSomething()
		if err != nil {
			ctx = wrapContext(ctx)
		}

		break
	}

	for {
		{
			ctx = wrapContext(ctx) // want "nested context in loop"
		}

		break
	}

	for {
		select {
		case <-ctx.Done():
			ctx = wrapContext(ctx) // want "nested context in loop"
		default:
		}

		break
	}
}

// this middleware could run on every request, bloating the request parameter level context and causing a memory leak
func badMiddleware(ctx context.Context) func() error {
	return func() error {
		ctx = wrapContext(ctx) // want "nested context in function literal"
		return doSomethingWithCtx(ctx)
	}
}

// this middleware is fine, as it doesn't modify the context of parent function
func okMiddleware(ctx context.Context) func() error {
	return func() error {
		ctx := wrapContext(ctx)
		return doSomethingWithCtx(ctx)
	}
}

// this middleware is fine, as it only modifies the context passed to it
func okMiddleware2(ctx context.Context) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ctx = wrapContext(ctx)
		return doSomethingWithCtx(ctx)
	}
}

func doSomethingWithCtx(ctx context.Context) error {
	return nil
}

func testCasesInit(t *testing.T) {
	cases := []struct {
		ctx context.Context
	}{
		{},
		{
			ctx: context.WithValue(context.Background(), "key", "value"),
		},
	}

	// This is fine because the first action is to re-assign to a known empty context.
	for _, tc := range cases {
		t.Run("some test", func(t *testing.T) {
			if tc.ctx == nil {
				tc.ctx = context.Background()
				tc.ctx = context.WithValue(tc.ctx, "key", "value")
			}
		})
	}
}

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func BeforeSuite(_ func()) interface{} {
	return nil
}

// This is fine because the first action is to re-assign to a known empty context.
var _ = BeforeSuite(func() {
	globalCtx = context.TODO()
	globalCtx, globalCancel = context.WithCancel(globalCtx)
})

func TestContext(t *testing.T) {
	tests := map[string]struct {
		ctx context.Context
	}{
		"WithContext":    {ctx: context.Background()},
		"WithoutContext": {},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Default to test context.
			if tc.ctx == nil {
				tc.ctx = t.Context()
			}
			fmt.Println(name, tc.ctx)
		})
	}
}
