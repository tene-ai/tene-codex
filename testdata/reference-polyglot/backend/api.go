package backend

import "context"

type Order struct {
	SKU      string
	Quantity int
}

func Create(ctx context.Context, o Order) error { return enqueue(ctx, o) }
func enqueue(context.Context, Order) error      { return nil }
