package store

import (
	"testing"
)

func TestCreateStore_UnsupportedType(t *testing.T) {
	_, err := CreateStore(StoreConfig{Type: "unknown"})
	if err == nil {
		t.Fatal("expected error for unsupported store type")
	}
	if err.Error() != "unsupported store type: unknown" {
		t.Errorf("err = %v", err)
	}
}

func TestDefaultStoreConfig(t *testing.T) {
	cfg := DefaultStoreConfig("json")
	if cfg.Type != "json" {
		t.Errorf("Type = %q, want json", cfg.Type)
	}
	if cfg.Settings == nil {
		t.Error("Settings should be initialized")
	}
}
