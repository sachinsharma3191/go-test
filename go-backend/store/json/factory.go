package json

import (
	"go-backend/store"
)

// Factory implements StoreFactory for JSON stores
type Factory struct{}

// CreateStore creates a new JSON store
func (f *Factory) CreateStore(config store.StoreConfig) (store.Store, error) {
	filePath, ok := config.Settings["file"].(string)
	if !ok {
		filePath = "data.json" // default
	}

	return NewStore(filePath)
}

// GetSupportedTypes returns the store types this factory supports
func (f *Factory) GetSupportedTypes() []string {
	return []string{"json"}
}

// Register the factory
func init() {
	store.RegisterFactory("json", &Factory{})
}
