package store

import (
	"context"
)

// Data represents any data model that can be stored
type Data interface {
	// Any model can implement this interface
}

// Store defines a minimal, model-agnostic data store interface
// Only provides Load and Save operations - all CRUD is handled by repositories
type Store interface {
	// Load loads the entire data from the store
	Load(ctx context.Context) (AppData, error)

	// Save saves the entire data to the store
	Save(ctx context.Context, data AppData) error

	// Health checks if the store is healthy
	Health(ctx context.Context) error

	// Close/cleanup resources
	Close() error
}
