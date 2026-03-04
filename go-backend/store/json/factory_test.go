package json

import (
	"testing"

	"go-backend/store"
)

func TestFactory_CreateStore(t *testing.T) {
	f := &Factory{}
	cfg := store.StoreConfig{
		Type:     "json",
		Settings: map[string]interface{}{"file": "factory_test.json"},
	}
	s, err := f.CreateStore(cfg)
	if err != nil {
		t.Fatalf("CreateStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestFactory_CreateStore_DefaultFile(t *testing.T) {
	f := &Factory{}
	cfg := store.StoreConfig{Type: "json", Settings: map[string]interface{}{}}
	s, err := f.CreateStore(cfg)
	if err != nil {
		t.Fatalf("CreateStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestFactory_GetSupportedTypes(t *testing.T) {
	f := &Factory{}
	types := f.GetSupportedTypes()
	if len(types) != 1 || types[0] != "json" {
		t.Errorf("GetSupportedTypes = %v", types)
	}
}
