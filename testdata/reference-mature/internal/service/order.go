package service

import (
	"context"
	"github.com/tene-ai/reference-mature/internal/store"
)

// Process preserves a legacy entry point while delegating persistence.
func Process(ctx context.Context, id string) error { return store.Save(ctx, id) }
