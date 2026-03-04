package store

import (
	"fmt"
)

// Builder helps create stores with different configurations
type Builder struct {
	config StoreConfig
}

// NewBuilder creates a new store builder
func NewBuilder(storeType string) *Builder {
	return &Builder{
		config: DefaultStoreConfig(storeType),
	}
}

// WithFile sets the file path for file-based stores
func (b *Builder) WithFile(path string) *Builder {
	b.config.Settings["file"] = path
	return b
}

// WithDatabase sets the database path for database stores
func (b *Builder) WithDatabase(path string) *Builder {
	b.config.Settings["db"] = path
	return b
}

// WithSetting adds a custom setting
func (b *Builder) WithSetting(key string, value interface{}) *Builder {
	b.config.Settings[key] = value
	return b
}

// Build creates the store
func (b *Builder) Build() (Store, error) {
	return CreateStore(b.config)
}

// QuickCreate provides a simple way to create stores
func QuickCreate(backend string, path string) (Store, error) {
	config := map[string]interface{}{}

	switch backend {
	case "json":
		config["file"] = path
	case "sqlite":
		config["db"] = path
	default:
		return nil, fmt.Errorf("unsupported backend: %s", backend)
	}

	builder := NewBuilder(backend)

	// Apply configuration
	for key, value := range config {
		switch key {
		case "file":
			builder = builder.WithFile(value.(string))
		case "db":
			builder = builder.WithDatabase(value.(string))
		default:
			builder = builder.WithSetting(key, value)
		}
	}

	return builder.Build()
}
