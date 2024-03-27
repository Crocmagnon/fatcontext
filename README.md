# foreshadow

`foreshadow` is a Go linter which detects un-shadowed contexts in loops.
They can lead to performance issues, as documented here: https://gabnotes.org/fat-contexts/

## Example

```go
package main

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
		ctx = context.WithValue(ctx, "key", i) // "context not shadowed in loop"
		_ = ctx
	}
}
```
