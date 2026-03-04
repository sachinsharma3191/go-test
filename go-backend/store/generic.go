package store

import (
	"fmt"
)

// StoreConfig holds configuration for store implementations
type StoreConfig struct {
	Type     string                 // "json", "sqlite", "postgres", etc.
	Settings map[string]interface{} // Implementation-specific settings
}

// StoreFactory creates stores based on configuration
type StoreFactory interface {
	CreateStore(config StoreConfig) (Store, error)
	GetSupportedTypes() []string
}

// Registry holds all registered store factories
var Registry = map[string]StoreFactory{}

// RegisterFactory registers a new store factory
func RegisterFactory(name string, factory StoreFactory) {
	Registry[name] = factory
}

// CreateStore creates a store from configuration
func CreateStore(config StoreConfig) (Store, error) {
	factory, exists := Registry[config.Type]
	if !exists {
		return nil, fmt.Errorf("unsupported store type: %s", config.Type)
	}

	return factory.CreateStore(config)
}

// DefaultStoreConfig returns a default configuration
func DefaultStoreConfig(storeType string) StoreConfig {
	return StoreConfig{
		Type:     storeType,
		Settings: make(map[string]interface{}),
	}
}
