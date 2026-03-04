// Integration tests for store - uses store_test package to avoid import cycle with store/json
package store_test

import (
	"os"
	"testing"

	"go-backend/store"

	_ "go-backend/store/json"
)

func TestBuilder_Build_JSON(t *testing.T) {
	f, _ := os.CreateTemp("", "builder_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, err := store.NewBuilder("json").WithFile(path).Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestQuickCreate_JSON(t *testing.T) {
	f, _ := os.CreateTemp("", "quick_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, err := store.QuickCreate("json", path)
	if err != nil {
		t.Fatalf("QuickCreate: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestCreateStore_JSON(t *testing.T) {
	f, _ := os.CreateTemp("", "createstore_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	cfg := store.StoreConfig{Type: "json", Settings: map[string]interface{}{"file": path}}
	s, err := store.CreateStore(cfg)
	if err != nil {
		t.Fatalf("CreateStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}
