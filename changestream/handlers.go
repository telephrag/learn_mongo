package changestream

import (
	"context"
	"fmt"
)

var handlerToEvent = map[string]func(ctx context.Context, cancel context.CancelFunc){
	"insert": func(ctx context.Context, cancel context.CancelFunc) {
		fmt.Println("insert handling...")
	},
	"delete": func(ctx context.Context, cancel context.CancelFunc) {
		fmt.Println("delete handling...")
	},
	"invalidate": func(ctx context.Context, cancel context.CancelFunc) {
		fmt.Println("invalidate handling...")
		cancel()
	},
}
