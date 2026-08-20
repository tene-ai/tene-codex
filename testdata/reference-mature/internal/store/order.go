package store

import (
	"context"
	"os"
)

func Save(_ context.Context, id string) error { return os.WriteFile("orders.db", []byte(id), 0600) }
