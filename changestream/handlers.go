package changestream

import (
	"context"
	"fmt"
)

var handlerToEvent = map[string]func(cancel context.CancelFunc){
	"insert": func(cancel context.CancelFunc) {
		fmt.Println("insert handling...")
	},
	"delete": func(cancel context.CancelFunc) {
		fmt.Println("delete handling...")
	},
	"invalidate": func(cancel context.CancelFunc) {
		fmt.Println("invalidate handling...")
		cancel()
	},
}
