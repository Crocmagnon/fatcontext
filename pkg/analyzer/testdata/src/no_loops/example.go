package no_loops

import (
	"context"
)

// Loop detection is disabled: this must NOT be reported.
func inLoop() {
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		ctx = context.WithValue(ctx, "key", i)
		_ = ctx
	}
}

// Function literal detection stays enabled: this MUST be reported.
func inFuncLit() {
	ctx := context.Background()

	f := func() {
		ctx = context.WithValue(ctx, "key", "val") // want "nested context in function literal"
		_ = ctx
	}
	f()
}
