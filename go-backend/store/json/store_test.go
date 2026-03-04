package json

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"go-backend/model"
	"go-backend/store"
)

func TestNewStore(t *testing.T) {
	f, _ := os.CreateTemp("", "newstore_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	dir, _ := os.MkdirTemp("", "jsondir")
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, "sub", "data.json")
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestStore_Load_NoFile(t *testing.T) {
	f, _ := os.CreateTemp("", "load_*.json")
	path := f.Name()
	os.Remove(path) // ensure file does not exist
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	data, err := s.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(data.Users) != 0 || len(data.Tasks) != 0 {
		t.Errorf("expected empty data: %+v", data)
	}
}

func TestStore_Load_EmptyFile(t *testing.T) {
	f, _ := os.CreateTemp("", "empty_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	data, err := s.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(data.Users) != 0 || len(data.Tasks) != 0 {
		t.Errorf("expected empty data")
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	f, _ := os.CreateTemp("", "save_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	data := store.AppData{
		Users: []model.User{{ID: 1, Name: "A", Email: "a@a.com", Role: "dev"}},
		Tasks: []model.Task{{ID: 1, Title: "T", Status: "pending", UserID: 1}},
	}
	err := s.Save(context.Background(), data)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := s.Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Users) != 1 || loaded.Users[0].Name != "A" {
		t.Errorf("loaded.Users = %v", loaded.Users)
	}
	if len(loaded.Tasks) != 1 || loaded.Tasks[0].Title != "T" {
		t.Errorf("loaded.Tasks = %v", loaded.Tasks)
	}
}

func TestStore_Health(t *testing.T) {
	f, _ := os.CreateTemp("", "health_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	err := s.Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
}

func TestStore_Close(t *testing.T) {
	f, _ := os.CreateTemp("", "close_*.json")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	err := s.Close()
	if err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "invalid_*.json")
	path := f.Name()
	os.WriteFile(path, []byte(`{invalid json}`), 0644)
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	s, _ := NewStore(path)
	_, err := s.Load(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestStore_NewStore_DirCreation(t *testing.T) {
	dir, _ := os.MkdirTemp("", "newstore")
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, "nested", "data.json")
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
}

func TestNewStore_MkdirAllFails(t *testing.T) {
	f, _ := os.CreateTemp("", "file_")
	path := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(path) })
	// path is a file; path/nested/data.json has parent path/nested whose parent is path (a file)
	// MkdirAll cannot create path/nested because path is a file
	storePath := filepath.Join(path, "nested", "data.json")
	_, err := NewStore(storePath)
	if err == nil {
		t.Fatal("expected NewStore to fail when parent is a file")
	}
}

func TestStore_Save_RenameFails(t *testing.T) {
	dir, _ := os.MkdirTemp("", "save_rename")
	t.Cleanup(func() { os.RemoveAll(dir) })
	// Use directory as filePath so Rename(tempFile, dir) fails (target is directory)
	s := &Store{filePath: dir}
	data := store.AppData{}
	err := s.Save(context.Background(), data)
	if err == nil {
		t.Fatal("expected Save to fail when filePath is a directory")
	}
}

func TestStore_Save_WriteFileFails(t *testing.T) {
	dir, _ := os.MkdirTemp("", "save_write")
	t.Cleanup(func() {
		os.Chmod(dir, 0755)
		os.RemoveAll(dir)
	})
	path := filepath.Join(dir, "data.json")
	s, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	// Make directory read-only so WriteFile to data.json.tmp fails
	if err := os.Chmod(dir, 0555); err != nil {
		t.Skipf("Cannot chmod on this system: %v", err)
	}
	data := store.AppData{}
	err = s.Save(context.Background(), data)
	if err == nil {
		t.Fatal("expected Save to fail in read-only directory")
	}
}